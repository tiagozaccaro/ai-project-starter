package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
)

// TestGeneratePlatformParity generates a project with all platform parity features
// and verifies that the expected files are generated for web, api, and mobile.
func TestGeneratePlatformParity(t *testing.T) {
	dest := t.TempDir()
	
	// Create a temporary config file with the required structure for platform parity.
	configContent := `
schemaVersion: 2
generatorVersion: "1.0.0"
profile: custom
projectName: platform-parity-test
applications:
  web: true
  api: true
  mobile: true
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  offlineStorage: true
  pushNotifications: true
  stateManagement: true
  deepLinking: true
  testing: true
metadata:
  description: "Platform parity test"
  license: "MIT"
  deepLinkScheme: "pptest"
  deepLinkUrl: "https://platform-parity-test.local"
`
	configPath := filepath.Join(dest, "project.config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	options := pipeline.GenerateOptions{
		ConfigPath:     configPath,
		Dest:           dest,
		NonInteractive: true,
	}
	err := pipeline.Generate(options)
	require.Nil(t, err, "platform parity generation should succeed")

	// Web application files
	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/hooks/use-offline-storage.ts"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/hooks/use-push-notifications.ts"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/hooks/use-state-management.ts"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/hooks/use-deep-linking.ts"))

	// API application files
	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/core/database.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/api/v1/sync.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/models/offline_sync.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/utils/conflict_resolution.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/core/security.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/core/config.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/api/v1/routes.py"))

	// Mobile application files
	assert.FileExists(t, filepath.Join(dest, "apps/mobile/pubspec.yaml"))
	assert.FileExists(t, filepath.Join(dest, "apps/mobile/lib/services/fcm_service.dart"))
	assert.FileExists(t, filepath.Join(dest, "apps/mobile/lib/state_management.dart"))
	assert.FileExists(t, filepath.Join(dest, "apps/mobile/lib/offline_storage.dart"))
	assert.FileExists(t, filepath.Join(dest, "apps/mobile/lib/deep_link_handler.dart"))

	// SC-005: When features are disabled, no output should be generated.
	// We test this by generating with only web and api (no mobile) and verifying mobile output is absent.
	// But note: platform parity features require both web and mobile, so they should be disabled when mobile is off.
	// We'll test this in a separate test case if needed.

	// For now, we just verify the positive case.
}