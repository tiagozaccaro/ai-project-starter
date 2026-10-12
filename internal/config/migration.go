package config

// Migration describes a single schema migration step.
type Migration struct {
	FromVersion int
	ToVersion   int
	Description string
}

// Migrations lists the ordered migration path from v1 to the current
// schema version.
func Migrations() []Migration {
	return []Migration{
		{FromVersion: 1, ToVersion: 2, Description: "Add generatorVersion, mobile.destinations, features.mcp; rename features.auth.method to features.auth"},
	}
}

// LatestSchemaVersion returns the current schema version.
func LatestSchemaVersion() int { return SchemaVersion }

// MigrationPath returns the ordered list of migrations needed to move
// from fromVersion to toVersion. It returns an error when no path
// exists (e.g., a downgrade or a gap).
func MigrationPath(fromVersion, toVersion int) ([]Migration, error) {
	if fromVersion == toVersion {
		return nil, nil
	}
	if fromVersion > toVersion {
		return nil, &MigrationError{
			Code:    "NO_MIGRATION_PATH",
			Message: "downgrade migrations are not supported",
		}
	}
	all := Migrations()
	var path []Migration
	current := fromVersion
	for {
		if current == toVersion {
			break
		}
		next := -1
		for _, m := range all {
			if m.FromVersion == current {
				next = m.ToVersion
				path = append(path, m)
				break
			}
		}
		if next == -1 {
			return nil, &MigrationError{
				Code:    "NO_MIGRATION_PATH",
				Message: "no migration path to target version",
			}
		}
		current = next
	}
	return path, nil
}

// MigrationError represents a migration failure.
type MigrationError struct {
	Code    string
	Message string
}

func (e *MigrationError) Error() string {
	return e.Code + ": " + e.Message
}

// MigrateV1ToV2 transforms a raw v1 configuration map into v2. It is
// pure and idempotent: running it on an already-v2 map is a no-op.
//
// v1 → v2 changes:
//   - adds generatorVersion (defaults to "1.0.0")
//   - adds features.mcp (default false)
//   - adds metadata.mobile.destinations when metadata.mobile exists
//   - renames features.auth.method → features.auth (list of providers)
func MigrateV1ToV2(raw map[string]interface{}) map[string]interface{} {
	if raw == nil {
		return raw
	}
	// Only migrate when schemaVersion is 1 (or absent, treated as v1).
	if v, ok := raw["schemaVersion"]; ok {
		if sv, ok := toInt(v); ok && sv >= 2 {
			return raw
		}
	}

	out := copyMap(raw)
	out["schemaVersion"] = 2

	if _, ok := out["generatorVersion"]; !ok {
		out["generatorVersion"] = GeneratorVersion
	}

	// features.mcp default false + auth.method rename.
	if feats, ok := out["features"].(map[string]interface{}); ok {
		if _, ok := feats["mcp"]; !ok {
			feats["mcp"] = false
		}
		if method, ok := feats["auth"].(map[string]interface{}); ok {
			if providers, ok := method["method"].(interface{}); ok {
				delete(method, "method")
				feats["auth"] = []interface{}{providers}
			}
		}
	}

	// metadata.mobile.destinations default.
	if meta, ok := out["metadata"].(map[string]interface{}); ok {
		if mobile, ok := meta["mobile"].(map[string]interface{}); ok {
			if _, ok := mobile["destinations"]; !ok {
				mobile["destinations"] = []interface{}{}
			}
		}
	}

	return out
}

// ApplyMigrations runs the ordered migrations from the detected version
// to the target version on a raw config map.
func ApplyMigrations(raw map[string]interface{}, targetVersion int) (map[string]interface{}, error) {
	current := SchemaVersion
	if v, ok := raw["schemaVersion"]; ok {
		if sv, ok := toInt(v); ok {
			current = sv
		}
	}
	path, err := MigrationPath(current, targetVersion)
	if err != nil {
		return nil, err
	}
	result := raw
	for _, m := range path {
		if m.FromVersion == 1 && m.ToVersion == 2 {
			result = MigrateV1ToV2(result)
		}
	}
	return result, nil
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

func copyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
