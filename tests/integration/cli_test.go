// Package integration tests the ai-project-starter CLI integration.
package integration

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/tiagozaccaro/ai-project-starter/internal/cli"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
)

// TestCLIGenerate tests the generate command through the CLI entry point.
func TestCLIGenerate(t *testing.T) {
	// This test exercises the CLI generate command to improve coverage
	// of the main function and Execute function.
	
	// We can't easily test the full CLI here without a subprocess,
	// but we can test that the commands are registered correctly.
	
	// Test that NewRootCommand creates a command with all subcommands
	root := cli.NewRootCommand()
	
	// Verify generate command exists
	cmd, _, err := root.Find([]string{"generate"})
	if err != nil {
		t.Errorf("generate command not found: %v", err)
	}
	if cmd.Name() != "generate" {
		t.Errorf("expected generate command, got %s", cmd.Name())
	}
	
	// Test other commands
	for _, name := range []string{"validate", "migrate-config", "list-profiles", "list-features", "upgrade", "completion"} {
		cmd, _, err := root.Find([]string{name})
		if err != nil {
			t.Errorf("%s command not found: %v", name, err)
		}
		if cmd.Name() != name {
			t.Errorf("expected %s command, got %s", name, cmd.Name())
		}
	}
}

func TestCLIValidationError(t *testing.T) {
	// Test that the CLI handles validation errors correctly
	// by running generate with an invalid config
	
	// Create a temporary invalid config file
	invalidConfig := `
schemaVersion: 2
generatorVersion: "1.0.0"
profile: custom
projectName: ""
applications:
  web: false
  api: false
  mobile: false
database:
  type: postgresql
  architecture: fastapi
features: {}
`
	
	// Write the invalid config to a temp file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	
	// Run generate with the invalid config
	options := pipeline.GenerateOptions{
		ConfigPath:     configPath,
		Dest:           filepath.Join(tmpDir, "out"),
		NonInteractive: true,
	}
	
	err := pipeline.Generate(options)
	if err == nil {
		t.Fatal("expected validation error for empty project name and no applications")
	}
	
	var cliErr *pipeline.CLIError
	if !errors.As(err, &cliErr) {
		t.Fatalf("expected CLIError, got %T", err)
	}
	
	if cliErr.Code != pipeline.ExitValidationError {
		t.Errorf("expected ExitValidationError (%d), got %d", pipeline.ExitValidationError, cliErr.Code)
	}
}

func TestCLIDryRun(t *testing.T) {
	// Test dry-run path
	validConfig := `
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: test-project
applications:
  web: true
  api: true
  mobile: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth:
    - google-oidc
  dashboard: true
  login: true
  testing: true
  github: true
  vscode: true
  speckit: true
metadata:
  mobile:
    bundleId: "com.example.testproject"
    destinations:
      - id: "home"
        label: "Home"
        icon: "home"
        route: "/home"
      - id: "dashboard"
        label: "Dashboard"
        icon: "dashboard"
        route: "/dashboard"
      - id: "profile"
        label: "Profile"
        icon: "person"
        route: "/profile"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	
	options := pipeline.GenerateOptions{
		ConfigPath:     configPath,
		Dest:           filepath.Join(tmpDir, "out"),
		NonInteractive: true,
		DryRun:         true,
	}
	
	err := pipeline.Generate(options)
	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}
}

func TestCLIJSONOutput(t *testing.T) {
	// Test JSON output path
	validConfig := `
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: test-project
applications:
  web: true
  api: true
  mobile: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth:
    - google-oidc
  dashboard: true
  login: true
  testing: true
  github: true
  vscode: true
  speckit: true
metadata:
  mobile:
    bundleId: "com.example.testproject"
    destinations:
      - id: "home"
        label: "Home"
        icon: "home"
        route: "/home"
      - id: "dashboard"
        label: "Dashboard"
        icon: "dashboard"
        route: "/dashboard"
      - id: "profile"
        label: "Profile"
        icon: "person"
        route: "/profile"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	
	options := pipeline.GenerateOptions{
		ConfigPath:     configPath,
		Dest:           filepath.Join(tmpDir, "out"),
		NonInteractive: true,
		Output:         "json",
	}
	
	err := pipeline.Generate(options)
	if err != nil {
		t.Fatalf("json output failed: %v", err)
	}
}

func TestCLIForceFlag(t *testing.T) {
	// Test force flag path
	validConfig := `
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: test-project
applications:
  web: true
  api: true
  mobile: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth:
    - google-oidc
  dashboard: true
  login: true
  testing: true
  github: true
  vscode: true
  speckit: true
metadata:
  mobile:
    bundleId: "com.example.testproject"
    destinations:
      - id: "home"
        label: "Home"
        icon: "home"
        route: "/home"
      - id: "dashboard"
        label: "Dashboard"
        icon: "dashboard"
        route: "/dashboard"
      - id: "profile"
        label: "Profile"
        icon: "person"
        route: "/profile"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	
	// First generate to create the destination
	options := pipeline.GenerateOptions{
		ConfigPath:     configPath,
		Dest:           filepath.Join(tmpDir, "out"),
		NonInteractive: true,
	}
	err := pipeline.Generate(options)
	if err != nil {
		t.Fatalf("first generate failed: %v", err)
	}
	
	// Second generate with force should work
	options.Force = true
	err = pipeline.Generate(options)
	if err != nil {
		t.Fatalf("force generate failed: %v", err)
	}
}