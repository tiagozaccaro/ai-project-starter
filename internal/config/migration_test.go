package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	migs := Migrations()
	require.NotEmpty(t, migs)
	assert.Equal(t, 1, migs[0].FromVersion)
	assert.Equal(t, 2, migs[0].ToVersion)
	assert.NotEmpty(t, migs[0].Description)
}

func TestLatestSchemaVersion(t *testing.T) {
	assert.Equal(t, SchemaVersion, LatestSchemaVersion())
}

func TestMigrationPathSameVersion(t *testing.T) {
	path, err := MigrationPath(2, 2)
	require.NoError(t, err)
	assert.Nil(t, path)
}

func TestMigrationPathV1ToV2(t *testing.T) {
	path, err := MigrationPath(1, 2)
	require.NoError(t, err)
	require.Len(t, path, 1)
	assert.Equal(t, 1, path[0].FromVersion)
	assert.Equal(t, 2, path[0].ToVersion)
}

func TestMigrationPathDowngrade(t *testing.T) {
	_, err := MigrationPath(2, 1)
	require.Error(t, err)
	var mErr *MigrationError
	require.ErrorAs(t, err, &mErr)
	assert.Equal(t, "NO_MIGRATION_PATH", mErr.Code)
	assert.Contains(t, mErr.Message, "downgrade")
}

func TestMigrationPathGap(t *testing.T) {
	_, err := MigrationPath(1, 5)
	require.Error(t, err)
	var mErr *MigrationError
	require.ErrorAs(t, err, &mErr)
	assert.Equal(t, "NO_MIGRATION_PATH", mErr.Code)
}

func TestMigrationError(t *testing.T) {
	err := &MigrationError{Code: "X", Message: "msg"}
	assert.Equal(t, "X: msg", err.Error())
}

func TestMigrateV1ToV2Nil(t *testing.T) {
	assert.Nil(t, MigrateV1ToV2(nil))
}

func TestMigrateV1ToV2Idempotent(t *testing.T) {
	v2 := map[string]interface{}{
		"schemaVersion":    2,
		"generatorVersion": "1.0.0",
	}
	out := MigrateV1ToV2(v2)
	assert.Equal(t, 2, out["schemaVersion"])
	// Already-v2 maps are returned unchanged (no-op).
	assert.Equal(t, v2, out)
}

func TestMigrateV1ToV2AddsDefaults(t *testing.T) {
	v1 := map[string]interface{}{
		"schemaVersion": 1,
		"projectName":   "app",
		"features":      map[string]interface{}{},
	}
	out := MigrateV1ToV2(v1)

	assert.Equal(t, 2, out["schemaVersion"])
	assert.Equal(t, GeneratorVersion, out["generatorVersion"])

	feats := out["features"].(map[string]interface{})
	assert.Equal(t, false, feats["mcp"])
}

func TestMigrateV1ToV2RenamesAuthMethod(t *testing.T) {
	v1 := map[string]interface{}{
		"schemaVersion": 1,
		"features": map[string]interface{}{
			"auth": map[string]interface{}{
				"method": "google-oidc",
			},
		},
	}
	out := MigrateV1ToV2(v1)

	feats := out["features"].(map[string]interface{})
	auth, ok := feats["auth"].([]interface{})
	require.True(t, ok, "auth should become a list")
	require.Len(t, auth, 1)
	assert.Equal(t, "google-oidc", auth[0])
}

func TestMigrateV1ToV2MobileDestinations(t *testing.T) {
	v1 := map[string]interface{}{
		"schemaVersion": 1,
		"metadata": map[string]interface{}{
			"mobile": map[string]interface{}{
				"bundleId": "com.example.app",
			},
		},
	}
	out := MigrateV1ToV2(v1)

	meta := out["metadata"].(map[string]interface{})
	mobile := meta["mobile"].(map[string]interface{})
	dests, ok := mobile["destinations"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, dests)
}

func TestMigrateV1ToV2PreservesExistingFields(t *testing.T) {
	v1 := map[string]interface{}{
		"schemaVersion":    1,
		"generatorVersion": "0.9.0",
		"features": map[string]interface{}{
			"mcp": true,
		},
	}
	out := MigrateV1ToV2(v1)
	// Existing generatorVersion and mcp are preserved.
	assert.Equal(t, "0.9.0", out["generatorVersion"])
	feats := out["features"].(map[string]interface{})
	assert.Equal(t, true, feats["mcp"])
}

func TestApplyMigrations(t *testing.T) {
	v1 := map[string]interface{}{
		"schemaVersion": 1,
		"projectName":   "app",
	}
	out, err := ApplyMigrations(v1, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, out["schemaVersion"])
}

func TestApplyMigrationsSameVersion(t *testing.T) {
	v2 := map[string]interface{}{"schemaVersion": 2}
	out, err := ApplyMigrations(v2, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, out["schemaVersion"])
}

func TestApplyMigrationsDowngrade(t *testing.T) {
	_, err := ApplyMigrations(map[string]interface{}{"schemaVersion": 2}, 1)
	require.Error(t, err)
}

func TestApplyMigrationsDetectsVersion(t *testing.T) {
	// Absent schemaVersion defaults to the current version,
	// so no migration path is needed.
	raw := map[string]interface{}{"projectName": "x"}
	out, err := ApplyMigrations(raw, 2)
	require.NoError(t, err)
	assert.Equal(t, raw, out)
}

func TestToInt(t *testing.T) {
	assert.Equal(t, 2, mustToInt(t, 2))
	assert.Equal(t, 2, mustToInt(t, int64(2)))
	assert.Equal(t, 2, mustToInt(t, float64(2)))
	_, ok := toInt("two")
	assert.False(t, ok)
}

func mustToInt(t *testing.T, v interface{}) int {
	t.Helper()
	n, ok := toInt(v)
	require.True(t, ok)
	return n
}

func TestCopyMap(t *testing.T) {
	src := map[string]interface{}{"a": 1, "b": "two"}
	dst := copyMap(src)
	assert.Equal(t, src, dst)
	// Independent copy.
	dst["a"] = 99
	assert.Equal(t, 1, src["a"])
}
