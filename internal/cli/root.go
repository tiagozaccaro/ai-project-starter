// Package cli defines the project-forge command-line
// interface root command.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tiagozaccaro/project-forge/internal/cli/commands"
	"github.com/tiagozaccaro/project-forge/internal/pipeline"
)

// Version is the generator version, set at build time.
var Version = "1.0.0"

// osExit is the function called when the CLI exits. Defaults to os.Exit
// but can be overridden in tests.
var OSExit = os.Exit

// NewRootCommand builds the root cobra command with all
// subcommands registered.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "project-forge",
		Short:         "Generate consistent, maintainable, AI-friendly project repositories",
		Long:          "project-forge generates configured software repositories (Next.js, FastAPI, Flutter) from a validated project.config.yaml.",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().BoolP("verbose", "v", false, "increase verbosity")
	root.PersistentFlags().BoolP("quiet", "q", false, "suppress non-error output")
	root.PersistentFlags().Bool("no-color", false, "disable colored output")

	root.AddCommand(commands.NewGenerateCommand())
	root.AddCommand(commands.NewValidateCommand())
	root.AddCommand(commands.NewMigrateConfigCommand())
	root.AddCommand(commands.NewListProfilesCommand())
	root.AddCommand(commands.NewListFeaturesCommand())
	root.AddCommand(commands.NewUpgradeCommand())
	root.AddCommand(commands.NewCompletionCommand())

	return root
}

// Execute runs the root command.
func Execute() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		var cliErr *pipeline.CLIError
		if ok := asCLIError(err, &cliErr); ok {
			fmt.Fprintln(os.Stderr, "Error:", cliErr.Message)
			OSExit(cliErr.Code)
			return
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		OSExit(1)
	}
}

func asCLIError(err error, target **pipeline.CLIError) bool {
	for err != nil {
		if e, ok := err.(*pipeline.CLIError); ok {
			*target = e
			return true
		}
		type unwrapper interface{ Unwrap() error }
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
			continue
		}
		break
	}
	return false
}
