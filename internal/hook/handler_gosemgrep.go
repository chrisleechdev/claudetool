package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/laurenkt/claudetool/internal/hook/semgreprules"
)

func init() {
	Register("go-semgrep", handleGoSemgrep)
}

// goSemgrepTimeout bounds the scan so a slow semgrep can't stall the agent.
const goSemgrepTimeout = 20 * time.Second

// lineRange is an inclusive span of 1-based line numbers.
type lineRange struct {
	start, end int
}

func handleGoSemgrep(in *Input) (*Output, error) {
	filePath := extractFilePath(in)
	if !isScannedGoFile(filePath) {
		return nil, nil
	}

	findings, err := scanWithEmbeddedRules(filePath)
	if err != nil {
		// A lint tool failure shouldn't block the edit
		return nil, nil
	}

	findings = onChangedLines(findings, changedLineRanges(in, filePath))
	if len(findings) == 0 {
		return nil, nil
	}

	report := fmt.Sprintf("semgrep found issues in %s:\n%s", filePath, formatGoSemgrepFindings(findings))
	if !hasBlockingFinding(findings) {
		return &Output{
			HookSpecificOutput: &HookSpecificOutput{
				HookEventName:     in.HookEventName,
				AdditionalContext: "Advisory (not blocking): " + report,
			},
		}, nil
	}
	return &Output{Decision: "block", Reason: report}, nil
}

func isScannedGoFile(filePath string) bool {
	if !strings.HasSuffix(filePath, ".go") || strings.HasSuffix(filePath, ".pb.go") {
		return false
	}
	return !strings.Contains("/"+filepath.ToSlash(filePath), "/vendor/")
}

func scanWithEmbeddedRules(filePath string) ([]semgrepFinding, error) {
	rulesDir, err := os.MkdirTemp("", "go-semgrep-rules-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(rulesDir)

	if err := writeEmbeddedRules(rulesDir); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), goSemgrepTimeout)
	defer cancel()

	// One spawn for all rules: per-spawn cost dominates per-rule cost.
	data, err := runSemgrepConfig(ctx, rulesDir, filePath, "--disable-version-check", "--metrics=off")
	if err != nil {
		return nil, err
	}
	return parseSemgrepFindings(data)
}

func writeEmbeddedRules(dir string) error {
	entries, err := fs.ReadDir(semgreprules.FS, ".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		rule, err := semgreprules.FS.ReadFile(entry.Name())
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, entry.Name()), rule, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// changedLineRanges returns the lines of the file on disk that the tool call wrote.
func changedLineRanges(in *Input, filePath string) []lineRange {
	var newStrings []string
	switch in.ToolName {
	case "Edit":
		var e EditInput
		if err := json.Unmarshal(in.ToolInput, &e); err != nil {
			return nil
		}
		newStrings = []string{e.NewString}
	case "MultiEdit":
		var m MultiEditInput
		if err := json.Unmarshal(in.ToolInput, &m); err != nil {
			return nil
		}
		for _, e := range m.Edits {
			newStrings = append(newStrings, e.NewString)
		}
	default:
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	var ranges []lineRange
	for _, newString := range newStrings {
		found := occurrenceLineRanges(string(content), newString)
		if len(found) == 0 {
			return nil
		}
		ranges = append(ranges, found...)
	}
	return ranges
}

func occurrenceLineRanges(content, text string) []lineRange {
	if text == "" {
		return nil
	}

	extraLines := strings.Count(strings.TrimSuffix(text, "\n"), "\n")
	var ranges []lineRange
	offset := 0
	for {
		i := strings.Index(content[offset:], text)
		if i < 0 {
			return ranges
		}
		start := offset + i
		startLine := 1 + strings.Count(content[:start], "\n")
		ranges = append(ranges, lineRange{start: startLine, end: startLine + extraLines})
		offset = start + len(text)
	}
}

// onChangedLines keeps the findings overlapping a changed range. No ranges
// means the whole file changed (a Write, or an Edit whose new text can't be
// located, including deletions).
func onChangedLines(findings []semgrepFinding, ranges []lineRange) []semgrepFinding {
	if len(ranges) == 0 {
		return findings
	}

	var kept []semgrepFinding
	for _, finding := range findings {
		for _, r := range ranges {
			if finding.StartLine <= r.end && finding.EndLine >= r.start {
				kept = append(kept, finding)
				break
			}
		}
	}
	return kept
}

func hasBlockingFinding(findings []semgrepFinding) bool {
	for _, finding := range findings {
		if finding.Severity == "ERROR" || finding.Severity == "WARNING" {
			return true
		}
	}
	return false
}

func formatGoSemgrepFindings(findings []semgrepFinding) string {
	var b strings.Builder
	for _, finding := range findings {
		fmt.Fprintf(&b, "  %s:%d: [%s] %s\n", finding.Path, finding.StartLine, ruleID(finding.CheckID), strings.TrimSpace(finding.Message))
		if finding.Fix != "" {
			fmt.Fprintf(&b, "    Suggested fix: %s\n", finding.Fix)
		}
	}
	return b.String()
}

// ruleID strips the path prefix semgrep prepends to check_id for rules loaded
// from a config directory.
func ruleID(checkID string) string {
	return checkID[strings.LastIndex(checkID, ".")+1:]
}
