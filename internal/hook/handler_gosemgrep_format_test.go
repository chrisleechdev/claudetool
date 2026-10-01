package hook

import "testing"

func TestFormatGoSemgrepFindings(t *testing.T) {
	findings := []semgrepFinding{
		{CheckID: "rules.first", Path: "a.go", StartLine: 2, Message: " first message \n", Fix: "replacement"},
		{CheckID: "rules.second", Path: "b.go", StartLine: 3, Message: "second message"},
		{CheckID: "rules.first", Path: "a.go", StartLine: 4, Message: " first message \n"},
	}
	want := "  a.go:2: [first]\n    Suggested fix: replacement\n  b.go:3: [second]\n  a.go:4: [first]\n\nRules:\n  [first] first message\n  [second] second message\n"
	if got := formatGoSemgrepFindings(findings); got != want {
		t.Errorf("report = %q, want %q", got, want)
	}
}
