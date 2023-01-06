package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/sunbeamlauncher/sunbeam/app"
	"github.com/sunbeamlauncher/sunbeam/tui"
	"github.com/sunbeamlauncher/sunbeam/utils"
)

func NewCmdExtension(api app.Api) *cobra.Command {
	extensionCommand := &cobra.Command{
		Use:     "extension",
		Aliases: []string{"extensions", "ext"},
		Short:   "Manage sunbeam extensions",
	}

	extensionArgs := make([]string, 0, len(api.Extensions))
	for _, extension := range api.Extensions {
		extensionArgs = append(extensionArgs, extension.Name)
	}

	extensionCommand.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:   "install <name> <repository>",
			Short: "Install a sunbeam extension from a git repository",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				extensionName := args[0]
				extensionRoot := args[1]
				if _, err := os.Stat(extensionRoot); err == nil {
					extensionRoot, err = filepath.Abs(extensionRoot)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Failed to get absolute path for extension root: %s", err)
						os.Exit(1)
					}

					if _, err = os.Stat(path.Join(extensionRoot, "sunbeam.yml")); os.IsNotExist(err) {
						return fmt.Errorf("current directory is not a sunbeam extension")
					}

					symlinkTarget := path.Join(api.ExtensionRoot, extensionName)

					if err := os.Symlink(extensionRoot, symlinkTarget); err != nil {
						fmt.Fprintln(os.Stderr, "Failed to create symlink", err)
						os.Exit(1)
					}

					fmt.Printf("Installed extension %s", extensionName)
					return nil
				}

				repo, err := utils.ParseWithHost(args[1], "github.com")
				if err != nil {
					return err
				}

				tmpDir, err := os.MkdirTemp(os.TempDir(), "sunbeam")
				if err != nil {
					return err
				}

				err = utils.GitClone(repo.Url(), tmpDir)
				if err != nil {
					return err
				}

				manifestPath := path.Join(tmpDir, "sunbeam.yml")
				if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
					return fmt.Errorf("extension %s does not have a sunbeam.yml manifest", repo.Name)
				}

				manifest, err := app.ParseManifest(extensionName, manifestPath)
				if err != nil {
					return err
				}

				if manifest.PostInstall != "" {
					cmd := exec.Command("sh", "-c", manifest.PostInstall)
					cmd.Dir = tmpDir
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					if err != nil {
						return err
					}
				}

				target := path.Join(api.ExtensionRoot, repo.Host, repo.Owner, repo.Name)
				os.MkdirAll(path.Dir(target), 0755)
				if err := os.Rename(tmpDir, target); err != nil {
					return err
				}

				fmt.Printf("Installed extension %s", args[0])
				return nil
			},
		}
	}())

	extensionCommand.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:       "remove",
			ValidArgs: extensionArgs,
			Short:     "Remove an installed extension",
			RunE: func(cmd *cobra.Command, args []string) error {
				extensionPath := path.Join(api.ExtensionRoot, args[0])
				if _, err := os.Stat(extensionPath); os.IsNotExist(err) {
					fmt.Fprintln(os.Stderr, "Extension not found")
					os.Exit(1)
				}

				if err := os.RemoveAll(extensionPath); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to remove extension")
					os.Exit(1)
				}

				fmt.Println("Removed extension", args[0])
				return nil
			},
		}
	}())

	extensionCommand.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:       "upgrade",
			Short:     "Upgrade installed extension",
			Args:      cobra.ExactArgs(1),
			ValidArgs: extensionArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				extensionDir := path.Join(api.ExtensionRoot, args[0])
				if _, err := os.Stat(extensionDir); os.IsNotExist(err) {
					fmt.Fprintln(os.Stderr, "Extension not found")
					os.Exit(1)
				}

				gc := utils.NewGitClient(extensionDir)

				currentVersion := gc.GetCurrentVersion()
				latestVersion, err := gc.GetLatestVersion()
				if err != nil {
					return err
				}

				if currentVersion == latestVersion {
					fmt.Printf("Extension %s is already up to date", args[0])
					return nil
				}

				return gc.Pull()
			},
		}
	}())

	extensionCommand.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:     "list",
			Short:   "List installed extensions",
			Aliases: []string{"ls"},
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				rows := make([][]string, 0)
				for _, extension := range api.Extensions {
					rows = append(rows, []string{extension.Name})
				}

				writer := tablewriter.NewWriter(os.Stdout)
				writer.SetBorder(false)
				writer.SetColumnSeparator(" ")
				writer.AppendBulk(rows)
				writer.Render()
			},
		}
	}())

	extensionCommand.AddCommand(func() *cobra.Command {
		command := cobra.Command{
			Use:   "browse",
			Short: "Enter a UI for browsing and installing extensions",
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				client := utils.NewGHClient("github.com")
				if err != nil {
					return err
				}
				res := struct {
					Items []struct {
						Name  string
						Owner struct {
							Login string
						}
						FullName    string `json:"full_name"`
						Description string
						HtmlURL     string `json:"html_url"`
					}
				}{}

				err = client.Get("search/repositories?q=topic:sunbeam-extension", &res)
				if err != nil {
					return err
				}

				extensionItems := make([]tui.ListItem, len(res.Items))
				for i, repo := range res.Items {
					item := tui.ListItem{
						Id:       strconv.Itoa(i),
						Title:    fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name),
						Subtitle: repo.Description,
					}

					if _, err := os.Stat(filepath.Join(api.ExtensionRoot, "github.com", repo.FullName)); err == nil {
						item.Accessories = []string{
							"Installed",
						}

						item.Actions = []tui.Action{
							{
								Title: "Remove Extension",
								Cmd:   tui.NewExecCmd(fmt.Sprintf("sunbeam extension remove %s", repo.Name)),
							},
							{
								Title: "Open in Browser",
								Cmd:   tui.NewOpenCmd(repo.HtmlURL),
							},
						}
					} else {
						item.Actions = []tui.Action{
							{
								Title: "Install Extension",
								Cmd:   tui.NewExecCmd(fmt.Sprintf("sunbeam extension install %s", repo.HtmlURL)),
							},
							{
								Title: "Open in Browser",
								Cmd:   tui.NewOpenCmd(repo.HtmlURL),
							},
						}
					}

					extensionItems[i] = item
				}

				list := tui.NewList("Browse Extensions")
				list.SetItems(extensionItems)
				root := tui.NewModel(list)

				return tui.Draw(root)
			},
		}
		return &command
	}())
	return extensionCommand
}
