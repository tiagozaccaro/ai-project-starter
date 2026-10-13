package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// update regenerates golden manifests when set:
//
//	go test ./tests/integration -run TestGolden -update
var update = flag.Bool("update", false, "update golden manifests")

// TestGolden compares the rendered template output of every
// built-in profile against a committed manifest of
// "sha256  path" lines. Any template, catalog condition, or
// config change that alters rendered bytes or the planned file
// set fails this test until the golden is regenerated with
// -update, catching template regressions (SC-005, SC-015).
func TestGolden(t *testing.T) {
	for _, profile := range profiles() {
		t.Run(profile, func(t *testing.T) {
			rc := loadResolved(t, profile+".yaml")
			actual := buildManifest(t, rc)

			path := goldenPath(profile)
			if *update {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("mkdir golden: %v", err)
				}
				if err := os.WriteFile(path, actual, 0o644); err != nil {
					t.Fatalf("write golden: %v", err)
				}
				t.Logf("updated golden manifest for %s (%d bytes)", profile, len(actual))
				return
			}

			expected, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden manifest (regenerate with -update): %v", err)
			}
			if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
				t.Fatalf("golden mismatch for %s (-update to regenerate):\n%s", profile, diff)
			}
		})
	}
}

// sha256Hex returns the hex-encoded SHA-256 of data.
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
