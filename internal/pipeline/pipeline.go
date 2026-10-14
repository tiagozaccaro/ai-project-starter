package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tiagozaccaro/ai-project-starter/internal/assets"
	"github.com/tiagozaccaro/ai-project-starter/internal/config"
	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
	"github.com/tiagozaccaro/ai-project-starter/internal/filesystem"
	"github.com/tiagozaccaro/ai-project-starter/internal/setup"
	"github.com/tiagozaccaro/ai-project-starter/internal/summary"
	"github.com/tiagozaccaro/ai-project-starter/internal/templates"
)

// GenerateOptions holds the options for the generate command.
type GenerateOptions struct {
	Profile        string
	Name           string
	Dest           string
	ConfigPath     string
	Layout         string
	Interactive    bool
	NonInteractive bool
	DryRun         bool
	Force          bool
	Output         string
}

// Generate is the end-to-end generation pipeline:
// load → validate → resolve → (dry-run | destination check →
// render → write → post-gen → summary).
func Generate(opts GenerateOptions) error {
	out := os.Stdout

	// 1. Load configuration.
	cfg, err := loadConfig(opts)
	if err != nil {
		return err
	}

	// 2. Validate the complete configuration (FR-003).
	if verr := cfg.Validate(); verr != nil {
		return ValidationError(verr.Error())
	}

	// 3. Resolve dependencies and auto-corrections.
	resolver := dependencies.NewResolver()
	rc, err := resolver.Resolve(cfg)
	if err != nil {
		return ValidationError(err.Error())
	}

	renderer := templates.NewRenderer(assets.FS)
	manifest := templates.TemplateManifests(rc)

	// 4. Dry-run: print the planned file tree and exit (FR-007).
	if opts.DryRun {
		planned, err := renderer.PlannedFiles(manifest, rc)
		if err != nil {
			return TemplateError(err.Error())
		}
		formatter := summary.NewFormatter(out)
		formatter.PrintDryRun(rc, planned)
		return nil
	}

	// 5. Destination safety check (FR-006, SC-008).
	destination := cfg.ResolveDestination()
	if err := filesystem.CheckDestination(destination, opts.Force); err != nil {
		var safetyErr *filesystem.SafetyError
		if isSafetyError(err, &safetyErr) {
			return DestinationError(safetyErr.Message)
		}
		return DestinationError(err.Error())
	}

	// When --force overwrites a non-empty destination, back it up first.
	if opts.Force {
		if empty, _ := filesystem.DestinationEmpty(destination); !empty {
			if _, err := filesystem.CreateBackup(destination); err != nil {
				return DestinationError(fmt.Sprintf("backup failed: %v", err))
			}
		}
	}

	// 6. Render templates.
	files, err := renderer.RenderAll(manifest, rc)
	if err != nil {
		return TemplateError(err.Error())
	}

	// 7. Write files atomically.
	writer := filesystem.NewWriter(destination, false)
	writeFiles := make([]filesystem.File, len(files))
	for i, f := range files {
		writeFiles[i] = filesystem.File{
			TargetPath: f.TargetPath,
			Content:    f.Content,
			Executable: f.Executable,
		}
	}
	result, err := writer.WriteFiles(writeFiles)
	if err != nil {
		return TemplateError(err.Error())
	}

	// 8. Post-generation setup.
	pg := setup.NewPostGeneration(destination)
	_ = pg.GitInit()
	_ = pg.HardenGitignore()
	violations := pg.ValidateIntegrity(rc)

	// 9. Summary (FR-005).
	formatter := summary.NewFormatter(out)
	if strings.EqualFold(opts.Output, "json") {
		if err := formatter.PrintJSON(rc, result.Written); err != nil {
			return TemplateError(err.Error())
		}
	} else {
		formatter.PrintSuccess(rc, result.Written)
	}

	if len(violations) > 0 {
		fmt.Fprintln(out, "\nIntegrity warnings:")
		for _, v := range violations {
			fmt.Fprintf(out, "  - %s\n", v)
		}
	}

	return nil
}

// loadConfig resolves the configuration source: --config file,
// --profile defaults, the AI_PROJECT_STARTER_CONFIG env var,
// or the interactive wizard. Env-var field overrides are
// applied last (flags > env > config file).
func loadConfig(opts GenerateOptions) (*config.ProjectConfig, error) {
	configPath := opts.ConfigPath
	if configPath == "" {
		if envPath, _ := config.EnvConfigSource(); envPath != "" {
			configPath = envPath
		}
	}

	switch {
	case configPath != "":
		result, err := config.LoadFromFile(configPath)
		if err != nil {
			return nil, ValidationError(err.Error())
		}
		cfg := result.Config
		if opts.Name != "" {
			cfg.ProjectName = opts.Name
		}
		if opts.Dest != "" {
			cfg.Destination = opts.Dest
		}
		if opts.Layout != "" {
			cfg.Layout = opts.Layout
		}
		config.ApplyEnvOverrides(cfg)
		return cfg, nil

	case opts.Profile != "":
		if !config.ProfileExists(opts.Profile) {
			valid := make([]string, 0)
			for _, p := range config.Profiles() {
				valid = append(valid, p.Name)
			}
			return nil, ValidationError(fmt.Sprintf(
				"unknown profile %q; valid profiles: %s",
				opts.Profile, strings.Join(valid, ", ")))
		}
		base, _ := config.DefaultProfileConfig(opts.Profile)
		cfg := base
		if opts.Name != "" {
			cfg.ProjectName = opts.Name
		} else {
			cfg.ProjectName = "my-project"
		}
		if opts.Dest != "" {
			cfg.Destination = opts.Dest
		}
		if opts.Layout != "" {
			cfg.Layout = opts.Layout
		}
		if opts.NonInteractive && cfg.ProjectName == "" {
			return nil, ValidationError("project name is required in non-interactive mode; use --name")
		}
		config.ApplyEnvOverrides(&cfg)
		return &cfg, nil

	case opts.Interactive || (opts.ConfigPath == "" && opts.Profile == "" && !opts.NonInteractive):
		return runWizard()

	default:
		return nil, ValidationError(
			"no configuration source; provide --profile, --config, or run interactively")
	}
}

// Validate runs configuration validation without generating.
func Validate(configPath string, strict bool, output string) error {
	result, err := config.LoadFromFile(configPath)
	if err != nil {
		return &CLIError{Code: 2, Message: err.Error()}
	}
	cfg := result.Config

	if verr := cfg.Validate(); verr != nil {
		if strings.EqualFold(output, "json") {
			printValidationJSON(os.Stdout, verr.Errors)
		} else {
			fmt.Fprintln(os.Stderr, verr.Error())
		}
		return ValidationError("configuration is invalid")
	}

	if strings.EqualFold(output, "json") {
		printValidationJSON(os.Stdout, nil)
	} else {
		fmt.Fprintln(os.Stdout, "Configuration is valid.")
	}

	// Auto-correction preview via the resolver.
	resolver := dependencies.NewResolver()
	rc, err := resolver.Resolve(cfg)
	if err != nil {
		return ValidationError(err.Error())
	}
	for _, ac := range rc.AutoCorrections {
		fmt.Fprintf(os.Stdout, "Auto-correction: %s %s -> %s (%s)\n", ac.Field, ac.From, ac.To, ac.Reason)
	}
	for _, w := range rc.Warnings {
		msg := fmt.Sprintf("[%s] %s", w.Code, w.Message)
		if strict {
			return ValidationError(msg)
		}
		fmt.Fprintf(os.Stdout, "Warning: %s\n", msg)
	}
	return nil
}

// MigrateConfig migrates a configuration file to a target schema version.
func MigrateConfig(configPath, outputPath string, noBackup bool, targetVersion int, dryRun bool) error {
	// Verify the config parses before migrating.
	if _, err := config.LoadFromFile(configPath); err != nil {
		return &CLIError{Code: 2, Message: err.Error()}
	}
	if targetVersion == 0 {
		targetVersion = config.LatestSchemaVersion()
	}

	// Load raw YAML for migration (preserves unknown fields).
	rawData, err := os.ReadFile(configPath)
	if err != nil {
		return &CLIError{Code: 2, Message: err.Error()}
	}
	var raw map[string]interface{}
	if err := yamlUnmarshal(rawData, &raw); err != nil {
		return &CLIError{Code: 2, Message: err.Error()}
	}

	current := config.SchemaVersion
	if v, ok := raw["schemaVersion"]; ok {
		if sv, ok := toInt(v); ok {
			current = sv
		}
	}

	if current == targetVersion {
		fmt.Fprintf(os.Stdout, "Config is already at schema version %d.\n", current)
		return nil
	}

	migrated, err := config.ApplyMigrations(raw, targetVersion)
	if err != nil {
		return &CLIError{Code: 3, Message: err.Error()}
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "Dry run: would migrate schemaVersion %d -> %d\n", current, targetVersion)
		return nil
	}

	out := outputPath
	if out == "" {
		out = configPath
		if !noBackup {
			backup := configPath + ".bak"
			if err := os.WriteFile(backup, rawData, 0o644); err != nil {
				return &CLIError{Code: 3, Message: err.Error()}
			}
			fmt.Fprintf(os.Stdout, "Backup created at %s\n", backup)
		}
	}

	data, err := yamlMarshal(migrated)
	if err != nil {
		return &CLIError{Code: 3, Message: err.Error()}
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return &CLIError{Code: 3, Message: err.Error()}
	}
	fmt.Fprintf(os.Stdout, "Migrated config to schema version %d at %s\n", targetVersion, out)
	return nil
}

// ListProfiles prints the available profiles.
func ListProfiles(output string) error {
	profiles := config.Profiles()
	switch strings.ToLower(output) {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(profiles)
	case "yaml":
		data, err := yamlMarshal(profiles)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	default:
		fmt.Printf("%-12s %s\n", "NAME", "DESCRIPTION")
		for _, p := range profiles {
			fmt.Printf("%-12s %s\n", p.Name, p.Description)
		}
		return nil
	}
}

// ListFeatures prints the configurable feature modules.
func ListFeatures(profile, output string) error {
	features := featureList(profile)
	switch strings.ToLower(output) {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(features)
	case "yaml":
		data, err := yamlMarshal(features)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	default:
		fmt.Printf("%-22s %-24s %s\n", "FEATURE", "REQUIRES", "CONFLICTS")
		for _, f := range features {
			fmt.Printf("%-22s %-24s %s\n", f.Name, f.Requires, f.Conflicts)
		}
		return nil
	}
}

// FeatureInfo describes a configurable feature module.
type FeatureInfo struct {
	Name      string `json:"name" yaml:"name"`
	Requires  string `json:"requires" yaml:"requires"`
	Conflicts string `json:"conflicts" yaml:"conflicts"`
}

func featureList(profile string) []FeatureInfo {
	features := []FeatureInfo{
		{"auth:google-oidc", "api", "auth:supabase, arch:supabase-first"},
		{"auth:supabase", "arch:supabase-first|hybrid", "auth:google-oidc"},
		{"dashboard", "web", "-"},
		{"login", "web", "-"},
		{"charts", "web", "-"},
		{"graphs", "web", "-"},
		{"ai", "api", "-"},
		{"workers", "api", "-"},
		{"observability", "web|api", "-"},
		{"docker", "-", "-"},
		{"vscode", "-", "-"},
		{"github", "-", "-"},
		{"speckit", "-", "-"},
		{"mcp", "-", "-"},
		{"testing", "-", "-"},
		{"database:migrations", "-", "db:mongodb (unless explicit)"},
	}
	if profile == "" {
		return features
	}
	base, ok := config.DefaultProfileConfig(profile)
	if !ok {
		return features
	}
	// Filter to features relevant to the profile's applications.
	var filtered []FeatureInfo
	for _, f := range features {
		filtered = append(filtered, f)
	}
	_ = base
	return filtered
}

func printValidationJSON(out io.Writer, errs []config.ValidationError) {
	type jsonErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
		Fix     string `json:"fix,omitempty"`
	}
	payload := struct {
		Valid  bool      `json:"valid"`
		Errors []jsonErr `json:"errors"`
	}{
		Valid:  len(errs) == 0,
		Errors: make([]jsonErr, 0, len(errs)),
	}
	for _, e := range errs {
		payload.Errors = append(payload.Errors, jsonErr{
			Code: e.Code, Message: e.Message, Field: e.Field, Fix: e.Fix,
		})
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

func isSafetyError(err error, target **filesystem.SafetyError) bool {
	if e, ok := err.(*filesystem.SafetyError); ok {
		*target = e
		return true
	}
	return false
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
}
