package cmd

import (
	"log"
	"strings"

	"github.com/pomdtr/sunbeam/api"
	"github.com/pomdtr/sunbeam/tui"
	"github.com/spf13/cobra"
)

type RunFlags struct {
	Params []string
}

var runFlags = RunFlags{}

func init() {
	runCmd.Flags().StringArrayVarP(&runFlags.Params, "param", "p", []string{}, "Parameters to pass to the script")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a sunbeam script",
	Run:   sunbeamRun,
	Args:  cobra.RangeArgs(1, 2),
}

func sunbeamRun(cmd *cobra.Command, args []string) {
	extensionName := args[0]

	manifest, ok := api.Sunbeam.Extensions[extensionName]
	if !ok {
		log.Fatalf("Extension %s not found", extensionName)
	}

	if len(args) < 2 {
		model := tui.RootList(manifest)
		err := tui.Draw(model, options)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	scriptName := args[1]

	script, ok := manifest.Pages[scriptName]
	if !ok {
		log.Fatalf("Page not found: %s", scriptName)
	}

	itemMap := make(map[string]api.FormItem)
	for _, formItem := range script.Inputs {
		itemMap[formItem.Name] = formItem
	}

	scriptParams := make(map[string]any)
	for _, param := range runFlags.Params {
		tokens := strings.SplitN(param, "=", 2)
		if len(tokens) != 2 {
			log.Fatalf("Invalid parameter: %s", param)
		}

		name, value := tokens[0], tokens[1]
		formItem, ok := itemMap[name]
		if !ok {
			log.Fatalf("Params %s does not exists in script %s", name, tokens[1])
		}
		switch formItem.Type {
		case "textfield":
			scriptParams[name] = value
		case "checkbox":
			scriptParams[name] = value == "true"
		}
	}

	container := tui.NewRunContainer(manifest, script, scriptParams)
	err := tui.Draw(container, options)
	if err != nil {
		log.Fatalf("could not run script: %v", err)
	}
}
