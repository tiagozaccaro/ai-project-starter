package dependencies

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tiagozaccaro/ai-project-starter/internal/config"
)

func TestAppDirMonorepo(t *testing.T) {
	assert.Equal(t, "apps/web", appDir(false, "apps/web"))
	assert.Equal(t, "apps/api", appDir(false, "apps/api"))
}

func TestAppDirFlat(t *testing.T) {
	assert.Equal(t, ".", appDir(true, "apps/web"))
	assert.Equal(t, ".", appDir(true, "apps/api"))
}

func TestResolvedDirsMonorepo(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{},
	}
	assert.Equal(t, "apps/web", rc.WebDir())
	assert.Equal(t, "apps/api", rc.APIDir())
	assert.Equal(t, "apps/mobile", rc.MobileDir())
}

func TestResolvedDirsFlat(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{},
		FlatLayout:    true,
	}
	assert.Equal(t, ".", rc.WebDir())
	assert.Equal(t, ".", rc.APIDir())
	assert.Equal(t, ".", rc.MobileDir())
}

func TestResolvedPathsMonorepo(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{},
	}
	assert.Equal(t, "apps/web/package.json", rc.WebPath("package.json"))
	assert.Equal(t, "apps/api/pyproject.toml", rc.APIPath("pyproject.toml"))
	assert.Equal(t, "apps/mobile/lib/main.dart", rc.MobilePath("lib/main.dart"))
}

func TestResolvedPathsFlat(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{},
		FlatLayout:    true,
	}
	assert.Equal(t, "package.json", rc.WebPath("package.json"))
	assert.Equal(t, "pyproject.toml", rc.APIPath("pyproject.toml"))
	assert.Equal(t, "lib/main.dart", rc.MobilePath("lib/main.dart"))
}

func TestResolvedString(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{ProjectName: "app"},
		GenerateWeb:   true,
		GenerateAPI:   true,
	}
	s := rc.String()
	assert.Contains(t, s, "app")
	assert.Contains(t, s, "web")
	assert.Contains(t, s, "api")
}

func TestResolvedStringMobile(t *testing.T) {
	rc := &ResolvedConfig{
		ProjectConfig:  &config.ProjectConfig{ProjectName: "app"},
		GenerateMobile: true,
	}
	assert.Contains(t, rc.String(), "mobile")
}
