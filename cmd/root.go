package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/cli/v2/pkg/findsh"
	"github.com/google/shlex"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"

	"github.com/pomdtr/sunbeam/internal"
	"github.com/pomdtr/sunbeam/types"
	"github.com/pomdtr/sunbeam/utils"
)

const (
	coreGroupID      = "core"
	extensionGroupID = "extension"
)

func Execute(version string) error {
	dataDir := filepath.Join(xdg.DataHome, "sunbeam")
	extensionRoot := filepath.Join(dataDir, "extensions")

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:          "sunbeam",
		Short:        "Command Line Launcher",
		Version:      version,
		SilenceUsage: true,
		Long: `Sunbeam is a command line launcher for your terminal, inspired by fzf and raycast.

See https://pomdtr.github.io/sunbeam for more information.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var input string
			if !isatty.IsTerminal(os.Stdin.Fd()) {
				return Run(internal.NewStaticGenerator(os.Stdin))
			}

			rootCommand, ok := os.LookupEnv("SUNBEAM_ROOT_CMD")
			if !ok {
				return cmd.Usage()
			}

			commandArgs, err := shlex.Split(rootCommand)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not parse default command: %s", err)
				return err
			}

			if len(commandArgs) == 0 {
				return cmd.Usage()
			}

			return Run(internal.NewCommandGenerator(&types.Command{
				Name:  commandArgs[0],
				Args:  commandArgs[1:],
				Input: input,
			}))
		},
	}

	rootCmd.Flags().StringArrayP("input", "i", nil, "input to pass to the action")
	rootCmd.Flags().String("query", "", "query to pass to the action")
	rootCmd.Flags().MarkHidden("input")
	rootCmd.Flags().MarkHidden("query")

	extensions, err := ListExtensions(extensionRoot)
	if err != nil {
		return fmt.Errorf("could not list extensions: %w", err)
	}

	rootCmd.AddGroup(
		&cobra.Group{ID: coreGroupID, Title: "Core Commands"},
		&cobra.Group{ID: extensionGroupID, Title: "Extension Commands"},
	)
	rootCmd.AddCommand(NewExtensionCmd(extensionRoot, extensions))
	rootCmd.AddCommand(NewQueryCmd())
	rootCmd.AddCommand(NewFetchCmd())
	rootCmd.AddCommand(NewReadCmd())
	rootCmd.AddCommand(NewTriggerCmd())
	rootCmd.AddCommand(NewValidateCmd())
	rootCmd.AddCommand(NewCmdRun(extensionRoot))
	rootCmd.AddCommand(NewInfoCmd(extensionRoot, version))

	rootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
	docCmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation for sunbeam",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			doc, err := buildDoc(rootCmd)
			if err != nil {
				return err
			}

			fmt.Println(doc)
			return nil
		},
	}
	rootCmd.AddCommand(docCmd)

	manCmd := &cobra.Command{
		Use:    "generate-man-pages [path]",
		Short:  "Generate Man Pages for sunbeam",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Title:   "MINE",
				Section: "3",
			}
			err := doc.GenManTree(rootCmd, header, args[0])
			if err != nil {
				return err
			}

			return nil
		},
	}
	rootCmd.AddCommand(manCmd)

	for extension, manifest := range extensions {
		rootCmd.AddCommand(NewExtensionExecCmd(extensionRoot, extension, manifest))
	}

	return rootCmd.Execute()
}

func NewExtensionExecCmd(extensionRoot string, extensionName string, manifest *ExtensionManifest) *cobra.Command {
	return &cobra.Command{
		Use:                extensionName,
		Short:              manifest.Description,
		DisableFlagParsing: true,
		GroupID:            extensionGroupID,

		RunE: func(cmd *cobra.Command, args []string) error {
			var input string
			if !isatty.IsTerminal(os.Stdin.Fd()) {
				inputBytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				input = string(inputBytes)
			}

			if manifest.Type == ExtentionTypeLocal {
				return runExtension(manifest.Entrypoint, args, input)
			}

			return runExtension(filepath.Join(extensionRoot, extensionName, manifest.Entrypoint), args, input)
		},
	}
}

func runExtension(extensionBin string, args []string, input string) error {
	var command types.Command
	if runtime.GOOS != "windows" {
		command = types.Command{
			Name: extensionBin,
			Args: args,
		}
		return Run(internal.NewCommandGenerator(&command))
	}

	shExe, err := findsh.Find()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("the `sh.exe` interpreter is required. Please install Git for Windows and try again")
		}
		return err
	}
	forwardArgs := append([]string{"-c", `command "$@"`, "--", extensionBin}, args...)

	command = types.Command{
		Name: shExe,
		Args: forwardArgs,
	}

	return Run(internal.NewCommandGenerator(&command))
}

func Run(generator internal.PageGenerator) error {
	if !isatty.IsTerminal(os.Stderr.Fd()) {
		output, err := generator()
		if err != nil {
			return fmt.Errorf("could not generate page: %s", err)
		}

		if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
			return fmt.Errorf("could not encode page: %s", err)
		}

		return nil

	}

	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		lipgloss.SetColorProfile(termenv.Ascii)
	} else {
		lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).Profile)
	}
	runner := internal.NewRunner(generator)

	return Draw(runner)
}

func Draw(page internal.Page) error {
	options := internal.SunbeamOptions{
		MaxHeight: utils.LookupInt("SUNBEAM_HEIGHT", 0),
		Padding:   utils.LookupInt("SUNBEAM_PADDING", 0),
	}
	paginator := internal.NewPaginator(page, options)

	var p *tea.Program
	if options.MaxHeight == 0 {
		p = tea.NewProgram(paginator, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	} else {
		p = tea.NewProgram(paginator, tea.WithOutput(os.Stderr))
	}

	m, err := p.Run()
	if err != nil {
		return err
	}

	paginator, ok := m.(*internal.Paginator)
	if !ok {
		return fmt.Errorf("could not cast model to paginator")
	}

	cmd := paginator.OutputCmd
	if cmd == nil {
		return nil
	}

	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func buildDoc(command *cobra.Command) (string, error) {
	if command.GroupID == extensionGroupID {
		return "", nil
	}

	var page strings.Builder
	err := doc.GenMarkdown(command, &page)
	if err != nil {
		return "", err
	}

	out := strings.Builder{}
	for _, line := range strings.Split(page.String(), "\n") {
		if strings.Contains(line, "SEE ALSO") {
			break
		}

		out.WriteString(line + "\n")
	}

	for _, child := range command.Commands() {
		childPage, err := buildDoc(child)
		if err != nil {
			return "", err
		}
		out.WriteString(childPage)
	}

	return out.String(), nil
}

func NewInfoCmd(extensionRoot string, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print information about sunbeam",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(func() (*types.Page, error) {
				return &types.Page{
					Title: "Info",
					Type:  types.ListPage,
					Items: []types.ListItem{
						{Title: "Version", Subtitle: version, Actions: []types.Action{
							types.NewCopyAction("Copy", version),
						}},
						{Title: "Extension Root", Subtitle: extensionRoot, Actions: []types.Action{
							types.NewCopyAction("Copy", extensionRoot),
						}},
					}}, nil
			})
		},
	}

}
