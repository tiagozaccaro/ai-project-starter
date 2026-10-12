package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFilesDryRun(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir, true)
	files := []File{
		{TargetPath: "a.txt", Content: []byte("A")},
		{TargetPath: "b.txt", Content: []byte("B")},
	}
	result, err := w.WriteFiles(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected dry-run result")
	}
	if result.FileCount != 2 {
		t.Errorf("expected 2 files, got %d", result.FileCount)
	}
	if len(result.Written) != 2 {
		t.Errorf("expected 2 written paths, got %d", len(result.Written))
	}
	if _, err := os.Stat(filepath.Join(dir, "a.txt")); !os.IsNotExist(err) {
		t.Error("dry-run should not write files")
	}
}

func TestWriteFilesReal(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "out")
	w := NewWriter(dir, false)
	files := []File{
		{TargetPath: "sub/a.txt", Content: []byte("hello")},
		{TargetPath: "sub/b.txt", Content: []byte("world")},
	}
	result, err := w.WriteFiles(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DryRun {
		t.Error("expected non-dry-run result")
	}
	if result.FileCount != 2 {
		t.Errorf("expected 2 files, got %d", result.FileCount)
	}
	data, err := os.ReadFile(filepath.Join(dir, "sub", "a.txt"))
	if err != nil {
		t.Fatalf("expected a.txt to exist: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected hello, got %q", string(data))
	}
}

func TestWriteFilesExecutableMode(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir, false)
	files := []File{
		{TargetPath: "run.sh", Content: []byte("#!/bin/sh"), Executable: true},
	}
	if _, err := w.WriteFiles(files); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, "run.sh"))
	if err != nil {
		t.Fatalf("expected run.sh to exist: %v", err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Errorf("expected 0755, got %v", info.Mode().Perm())
	}
}

func TestWriteFilesRegularMode(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir, false)
	files := []File{
		{TargetPath: "f.txt", Content: []byte("x")},
	}
	if _, err := w.WriteFiles(files); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, "f.txt"))
	if err != nil {
		t.Fatalf("expected f.txt to exist: %v", err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Errorf("expected 0644, got %v", info.Mode().Perm())
	}
}

func TestWriteFilesEmpty(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir, false)
	result, err := w.WriteFiles(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FileCount != 0 {
		t.Errorf("expected 0 files, got %d", result.FileCount)
	}
}

func TestListDestinationFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "b.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := ListDestinationFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(files), files)
	}
}
