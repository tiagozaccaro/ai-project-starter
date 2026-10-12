package commands

import (
	"github.com/spf13/cobra"

	pipeline "github.com/your-org/ai-project-starter/internal/pipeline"
)

// NewGenerateCommand creates the `generate` command.
func NewGenerateCommand() *cobra.Command {
	var opts pipeline.GenerateOptions

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from a profile, interactive config, or config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Generate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Profile, "profile", "p", "", "profile: fullstack, backend, frontend, ai-fullstack, custom")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "project name (kebab-case)")
	cmd.Flags().StringVarP(&opts.Dest, "dest", "d", "", "destination directory (default: ./{name})")
	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "path to project.config.yaml")
	cmd.Flags().StringVar(&opts.Layout, "layout", "", "project layout: monorepo, flat")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "enable interactive configuration wizard")
	cmd.Flags().BoolVar(&opts.NonInteractive, "non-interactive", false, "disable interactive prompts; fail if config incomplete")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "show file tree without writing")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "overwrite non-empty destination")
	cmd.Flags().StringVar(&opts.Output, "output", "", "output format: json, yaml")

	return cmd
}

// NewValidateCommand creates the `validate` command.
func NewValidateCommand() *cobra.Command {
	var configPath string
	var strict bool
	var output string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a project.config.yaml file without generating",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Validate(configPath, strict, output)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file (required)")
	cmd.Flags().BoolVar(&strict, "strict", false, "treat warnings as errors")
	cmd.Flags().StringVar(&output, "output", "", "output format: json, yaml")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

// NewMigrateConfigCommand creates the `migrate-config` command.
func NewMigrateConfigCommand() *cobra.Command {
	var configPath string
	var outputPath string
	var noBackup bool
	var targetVersion int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "migrate-config",
		Short: "Migrate a config file to the current schema version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.MigrateConfig(configPath, outputPath, noBackup, targetVersion, dryRun)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file (required)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path (default: overwrite in place with backup)")
	cmd.Flags().BoolVar(&noBackup, "no-backup", false, "skip creating .bak file")
	cmd.Flags().IntVar(&targetVersion, "target-version", 0, "target schema version (default: latest)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show diff without writing")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

// NewListProfilesCommand creates the `list-profiles` command.
func NewListProfilesCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "list-profiles",
		Short: "List available profiles with descriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.ListProfiles(output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "table", "output format: table, json, yaml")

	return cmd
}

// NewListFeaturesCommand creates the `list-features` command.
func NewListFeaturesCommand() *cobra.Command {
	var profile string
	var output string

	cmd := &cobra.Command{
		Use:   "list-features",
		Short: "List all configurable feature modules with dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.ListFeatures(profile, output)
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", "", "show features for a specific profile")
	cmd.Flags().StringVarP(&output, "output", "o", "table", "output format: table, json, yaml")

	return cmd
}
