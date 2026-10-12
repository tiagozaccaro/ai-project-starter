package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SafetyError is returned when destination protection
// prevents generation.
type SafetyError struct {
	Code        string
	Message     string
	Destination string
}

func (e *SafetyError) Error() string {
	return e.Message
}

// CheckDestination verifies that it is safe to generate
// into destination. A non-empty destination without
// force is rejected (exit code 2). An empty or missing
// destination is allowed.
func CheckDestination(destination string, force bool) error {
	empty, err := DestinationEmpty(destination)
	if err != nil {
		return &SafetyError{
			Code:        "DESTINATION_ERROR",
			Message:     fmt.Sprintf("cannot inspect destination %s: %v", destination, err),
			Destination: destination,
		}
	}
	if empty {
		return nil
	}
	if !force {
		return &SafetyError{
			Code:        "DESTINATION_NOT_EMPTY",
			Message:     fmt.Sprintf("destination %s is not empty; use --force to overwrite (a backup will be created)", destination),
			Destination: destination,
		}
	}
	return nil
}

// DetectConflicts reports which planned target paths
// already exist in the destination.
func DetectConflicts(destination string, planned []string) ([]string, error) {
	var conflicts []string
	for _, rel := range planned {
		target := filepath.Join(destination, rel)
		if _, err := os.Stat(target); err == nil {
			conflicts = append(conflicts, rel)
		}
	}
	return conflicts, nil
}

// CreateBackup renames the destination directory to a
// timestamped backup and returns the backup path.
func CreateBackup(destination string) (string, error) {
	info, err := os.Stat(destination)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("destination %s is not a directory", destination)
	}
	backup := fmt.Sprintf("%s.backup-%s", destination,
		time.Now().Format("20060102-150405"))
	if err := os.Rename(destination, backup); err != nil {
		return "", fmt.Errorf("create backup: %w", err)
	}
	return backup, nil
}
