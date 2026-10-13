package templates

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncsRegistered(t *testing.T) {
	funcs := Funcs()
	for _, name := range []string{
		"toSnake", "toKebab", "toPascal", "toCamel",
		"dartPackageName", "goModuleName", "dotenvEscape",
		"jsonEscape", "yamlEscape", "upper", "lower", "title",
	} {
		assert.Contains(t, funcs, name)
	}
}

func TestToSnake(t *testing.T) {
	assert.Equal(t, "my_app", toSnake("my-app"))
	assert.Equal(t, "my_app", toSnake("my app"))
	assert.Equal(t, "my_app", toSnake("myApp"))
	assert.Equal(t, "my_app", toSnake("MyApp"))
	assert.Equal(t, "already_snake", toSnake("already_snake"))
}

func TestToKebab(t *testing.T) {
	assert.Equal(t, "my-app", toKebab("my_app"))
	assert.Equal(t, "my-app", toKebab("my app"))
	assert.Equal(t, "my-app", toKebab("my.app"))
	assert.Equal(t, "a-b-c-d", toKebab("a_b.c d"))
}

func TestToPascal(t *testing.T) {
	assert.Equal(t, "MyApp", toPascal("my-app"))
	assert.Equal(t, "MyApp", toPascal("my_app"))
	assert.Equal(t, "MyApp", toPascal("my app"))
	assert.Equal(t, "ABC", toPascal("a-b-c"))
	assert.Equal(t, "", toPascal(""))
	// Consecutive separators produce empty segments that
	// are skipped.
	assert.Equal(t, "MyApp", toPascal("my--app"))
	assert.Equal(t, "MyApp", toPascal("-my-app-"))
}

func TestToCamel(t *testing.T) {
	assert.Equal(t, "myApp", toCamel("my-app"))
	assert.Equal(t, "myApp", toCamel("my_app"))
	assert.Equal(t, "", toCamel(""))
}

func TestDartPackageName(t *testing.T) {
	assert.Equal(t, "my_app", dartPackageName("my-app"))
	assert.Equal(t, "my_app", dartPackageName("MyApp"))
}

func TestGoModuleName(t *testing.T) {
	assert.Equal(t, "myapp", goModuleName("my-app"))
	assert.Equal(t, "myapp", goModuleName("myapp"))
}

func TestDotenvEscape(t *testing.T) {
	assert.Equal(t, "line1\\nline2", dotenvEscape("line1\nline2"))
	assert.Equal(t, `say \"hi\"`, dotenvEscape(`say "hi"`))
}

func TestJSONEscape(t *testing.T) {
	assert.Equal(t, `a\"b`, jsonEscape(`a"b`))
	assert.Equal(t, `a\\b`, jsonEscape(`a\b`))
	assert.Equal(t, "a\\nb", jsonEscape("a\nb"))
	assert.Equal(t, "a\\tb", jsonEscape("a\tb"))
	assert.Equal(t, "plain", jsonEscape("plain"))
}

func TestYAMLEscape(t *testing.T) {
	assert.Equal(t, `""`, yamlEscape(""))
	assert.Equal(t, "plain", yamlEscape("plain"))
	assert.Equal(t, `"has: colon"`, yamlEscape("has: colon"))
	assert.Equal(t, `"hash #"`, yamlEscape("hash #"))
	assert.Equal(t, `"trailing "`, yamlEscape("trailing "))
	assert.Equal(t, `"quote\"d"`, yamlEscape(`quote"d`))
}

func TestTitle(t *testing.T) {
	assert.Equal(t, "Hello", title("hello"))
	assert.Equal(t, "Hello", title("Hello"))
	assert.Equal(t, "", title(""))
}

func TestUpperLower(t *testing.T) {
	assert.Equal(t, "UP", Funcs()["upper"].(func(string) string)("up"))
	assert.Equal(t, "down", Funcs()["lower"].(func(string) string)("DOWN"))
}

func TestEvaluateTargetPath(t *testing.T) {
	rc := resolvedConfig()

	// Static target path is returned unchanged.
	got, err := EvaluateTargetPath("out/static.txt", rc)
	require.NoError(t, err)
	assert.Equal(t, "out/static.txt", got)

	// Templated target path is rendered.
	got, err = EvaluateTargetPath("{% .ProjectName %}.txt", rc)
	require.NoError(t, err)
	assert.Equal(t, "my-app.txt", got)
}

func TestEvaluateTargetPathInvalid(t *testing.T) {
	// A bare template that renders to a path is valid; an
	// unterminated tag is not.
	_, err := EvaluateTargetPath("{% .ProjectName %}", resolvedConfig())
	require.NoError(t, err)
	_, err = EvaluateTargetPath("{% if .GenerateWeb %}", resolvedConfig())
	assert.Error(t, err)
}

func TestCatalogFS(t *testing.T) {
	r := &Renderer{}
	assert.Nil(t, r.CatalogFS())
}

func TestRenderAllExecutable(t *testing.T) {
	fsys := fstest.MapFS{
		"script.tmpl": &fstest.MapFile{Data: []byte("#!/bin/sh")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "script.tmpl", TargetPath: "out/script.sh", Condition: ``, Executable: true},
	}
	r := &Renderer{FS: fsys}
	files, err := r.RenderAll(manifest, resolvedConfig())
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.True(t, files[0].Executable)
}
