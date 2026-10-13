package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerationTiming verifies full generation
// completes well within the 30-second budget (SC-001
// fullstack < 2min, SC-002 config-file < 30s).
func TestGenerationTiming(t *testing.T) {
	for _, profile := range profiles() {
		t.Run(profile, func(t *testing.T) {
			dest := t.TempDir()
			start := time.Now()
			cliErr := generateTo(t, profile+".yaml", dest)
			elapsed := time.Since(start)
			require.Nil(t, cliErr, "%s generation should succeed", profile)
			assert.Less(t, elapsed, 30*time.Second,
				"%s generation took %s, exceeding the 30s budget", profile, elapsed)
			t.Logf("%s generated in %s", profile, elapsed)
		})
	}
}
