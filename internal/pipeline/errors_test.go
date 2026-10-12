package pipeline

import "testing"

func TestCLIErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		wantCode int
		wantMsg  string
	}{
		{"validation", ValidationError("bad config"), ExitValidationError, "bad config"},
		{"destination", DestinationError("not empty"), ExitDestinationError, "not empty"},
		{"template", TemplateError("render failed"), ExitTemplateError, "render failed"},
		{"cancelled", CancelledError("user aborted"), ExitCancelled, "user aborted"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Code != tc.wantCode {
				t.Errorf("expected code %d, got %d", tc.wantCode, tc.err.Code)
			}
			if tc.err.Error() != tc.wantMsg {
				t.Errorf("expected message %q, got %q", tc.wantMsg, tc.err.Error())
			}
		})
	}
}

func TestExitCodeConstants(t *testing.T) {
	if ExitSuccess != 0 {
		t.Errorf("expected ExitSuccess 0, got %d", ExitSuccess)
	}
	if ExitValidationError != 1 {
		t.Errorf("expected ExitValidationError 1, got %d", ExitValidationError)
	}
	if ExitDestinationError != 2 {
		t.Errorf("expected ExitDestinationError 2, got %d", ExitDestinationError)
	}
	if ExitTemplateError != 3 {
		t.Errorf("expected ExitTemplateError 3, got %d", ExitTemplateError)
	}
	if ExitCancelled != 4 {
		t.Errorf("expected ExitCancelled 4, got %d", ExitCancelled)
	}
}

func TestCLIErrorIsDistinct(t *testing.T) {
	a := ValidationError("x")
	b := DestinationError("x")
	if a.Code == b.Code {
		t.Error("expected different exit codes for different error types")
	}
}
