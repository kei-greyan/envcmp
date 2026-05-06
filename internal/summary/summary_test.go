package summary_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/stats"
	"github.com/user/envcmp/internal/summary"
)

func baseStats(matched, missingLeft, missingRight, mismatched int) stats.Stats {
	return stats.Stats{
		Total:          matched + missingLeft + missingRight + mismatched,
		Matched:        matched,
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:     mismatched,
	}
}

func TestBuild_CleanResult(t *testing.T) {
	s := summary.Build(baseStats(5, 0, 0, 0))
	if !s.Clean {
		t.Error("expected Clean to be true")
	}
	if s.TotalKeys != 5 || s.Matched != 5 {
		t.Errorf("unexpected counts: %+v", s)
	}
}

func TestBuild_DirtyResult(t *testing.T) {
	s := summary.Build(baseStats(3, 1, 2, 1))
	if s.Clean {
		t.Error("expected Clean to be false")
	}
	if s.MissingLeft != 1 || s.MissingRight != 2 || s.Mismatched != 1 {
		t.Errorf("unexpected counts: %+v", s)
	}
}

func TestOneLiner_CleanMessage(t *testing.T) {
	s := summary.Build(baseStats(4, 0, 0, 0))
	line := s.OneLiner()
	if !strings.HasPrefix(line, "✓") {
		t.Errorf("expected clean prefix, got: %s", line)
	}
	if !strings.Contains(line, "4 keys matched") {
		t.Errorf("expected matched count in message, got: %s", line)
	}
}

func TestOneLiner_DirtyMessage(t *testing.T) {
	s := summary.Build(baseStats(2, 1, 1, 1))
	line := s.OneLiner()
	if !strings.HasPrefix(line, "✗") {
		t.Errorf("expected dirty prefix, got: %s", line)
	}
	for _, substr := range []string{"missing in left", "missing in right", "mismatched"} {
		if !strings.Contains(line, substr) {
			t.Errorf("expected %q in message, got: %s", substr, line)
		}
	}
}

func TestOneLiner_TotalKeyCount(t *testing.T) {
	s := summary.Build(baseStats(1, 2, 3, 4))
	line := s.OneLiner()
	if !strings.Contains(line, "10 total keys") {
		t.Errorf("expected total key count in message, got: %s", line)
	}
}
