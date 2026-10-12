package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDestinationEmptyMissing(t *testing.T) {
	empty, err := DestinationEmpty(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !empty {
		t.Error("expected missing destination to be empty")
	}
}

func TestDestinationEmptyDir(t *testing.T) {
	dir := t.TempDir()
	empty, err := DestinationEmpty(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !empty {
		t.Error("expected empty dir to be empty")
	}
}

func TestDestinationEmptyNonEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	empty, err := DestinationEmpty(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if empty {
		t.Error("expected non-empty dir to not be empty")
	}
}

func TestCheckDestinationEmpty(t *testing.T) {
	if err := CheckDestination(filepath.Join(t.TempDir(), "missing"), false); err != nil {
		t.Errorf("expected empty destination to be allowed, got %v", err)
	}
}

func TestCheckDestinationNonEmptyNoForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := CheckDestination(dir, false)
	if err == nil {
		t.Fatal("expected error for non-empty destination without force")
	}
	se, ok := err.(*SafetyError)
	if !ok {
		t.Fatalf("expected *SafetyError, got %T", err)
	}
	if se.Code != "DESTINATION_NOT_EMPTY" {
		t.Errorf("expected DESTINATION_NOT_EMPTY, got %q", se.Code)
	}
}

func TestCheckDestinationNonEmptyForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CheckDestination(dir, true); err != nil {
		t.Errorf("expected force to allow non-empty destination, got %v", err)
	}
}

func TestDetectConflicts(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	planned := []string{"existing.txt", "new.txt"}
	conflicts, err := DetectConflicts(dir, planned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 || conflicts[0] != "existing.txt" {
		t.Errorf("expected [existing.txt], got %v", conflicts)
	}
}

func TestDetectConflictsNone(t *testing.T) {
	conflicts, err := DetectConflicts(t.TempDir(), []string{"a.txt", "b.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", conflicts)
	}
}

func TestCreateBackup(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	backup, err := CreateBackup(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(backup); err != nil {
		t.Errorf("expected backup to exist: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("expected original destination to be renamed away")
	}
}

func TestCreateBackupNotDir(t *testing.T) {
	f := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateBackup(f); err == nil {
		t.Error("expected error backing up a non-directory")
	}
}
