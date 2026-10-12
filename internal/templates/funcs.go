package templates

import (
	"strings"
	"text/template"
)

// Funcs returns the custom template function map available to
// all project templates.
func Funcs() template.FuncMap {
	return template.FuncMap{
		"toSnake":         toSnake,
		"toKebab":         toKebab,
		"toPascal":        toPascal,
		"toCamel":         toCamel,
		"dartPackageName": dartPackageName,
		"goModuleName":    goModuleName,
		"dotenvEscape":    dotenvEscape,
		"jsonEscape":      jsonEscape,
		"yamlEscape":      yamlEscape,
		"upper":           strings.ToUpper,
		"lower":           strings.ToLower,
		"title":           title,
	}
}

// toSnake converts "kebab-case" or "PascalCase" to snake_case.
func toSnake(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	// Handle camelCase/PascalCase boundaries.
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + ('a' - 'A'))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// toKebab converts a string to kebab-case.
func toKebab(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "_", "-"), " ", "-"), ".", "-")
}

// toPascal converts "kebab-case" to PascalCase.
func toPascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

// toCamel converts "kebab-case" to camelCase.
func toCamel(s string) string {
	p := toPascal(s)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// dartPackageName converts a project name to a valid Dart
// package name (lowercase, underscores).
func dartPackageName(s string) string {
	return toSnake(s)
}

// goModuleName converts a project name to a valid Go module
// name (no hyphens).
func goModuleName(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

// dotenvEscape escapes a value for inclusion in a .env file.
func dotenvEscape(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// jsonEscape escapes a string for safe inclusion in JSON.
func jsonEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString("\\\"")
		case '\\':
			b.WriteString("\\\\")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// yamlEscape quotes a string if it needs escaping in YAML.
func yamlEscape(s string) string {
	if s == "" {
		return `""`
	}
	if strings.ContainsAny(s, ":#\n\"'{}[]&*!|>%@`,\t") || strings.HasSuffix(s, " ") {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}

func title(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
