package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// semgrepHandler returns a Handler that runs semgrep with the given YAML rule
// against files written/edited by Claude. Use as a PostToolUse hook.
func semgrepHandler(ruleYAML string) Handler {
	return func(in *Input) (*Output, error) {
		filePath := extractFilePath(in)
		if filePath == "" {
			return nil, nil
		}

		findings, err := runSemgrep(ruleYAML, filePath)
		if err != nil {
			// semgrep not installed or failed; don't block
			return nil, nil
		}
		if findings == "" {
			return nil, nil
		}

		return &Output{
			Decision: "block",
			Reason:   fmt.Sprintf("semgrep found issues in %s:\n%s", filePath, findings),
		}, nil
	}
}

func extractFilePath(in *Input) string {
	if len(in.ToolInput) == 0 {
		return ""
	}
	var obj struct {
		FilePath string `json:"file_path"`
	}
	if err := json.Unmarshal(in.ToolInput, &obj); err != nil {
		return ""
	}
	return obj.FilePath
}

type semgrepFinding struct {
	CheckID   string
	Path      string
	StartLine int
	EndLine   int
	Message   string
	Severity  string
	Fix       string
}

func runSemgrep(ruleYAML, filePath string) (string, error) {
	f, err := os.CreateTemp("", "semgrep-rule-*.yml")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(ruleYAML); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	data, err := runSemgrepConfig(context.Background(), f.Name(), filePath)
	if err != nil {
		return "", err
	}
	return parseSemgrepOutput(data), nil
}

// runSemgrepConfig scans filePath with the rules at configPath and returns semgrep's JSON output.
func runSemgrepConfig(ctx context.Context, configPath, filePath string, extraFlags ...string) ([]byte, error) {
	args := append([]string{"scan", "--quiet", "--json"}, extraFlags...)
	args = append(args, "--config", configPath, filePath)

	cmd := exec.CommandContext(ctx, "semgrep", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

func parseSemgrepOutput(data []byte) string {
	findings, err := parseSemgrepFindings(data)
	if err != nil {
		return string(data)
	}
	if len(findings) == 0 {
		return ""
	}

	var b strings.Builder
	for _, f := range findings {
		fmt.Fprintf(&b, "  %s:%d: [%s] %s\n", f.Path, f.StartLine, f.CheckID, f.Message)
	}
	return b.String()
}

func parseSemgrepFindings(data []byte) ([]semgrepFinding, error) {
	var result struct {
		Results []struct {
			CheckID string `json:"check_id"`
			Path    string `json:"path"`
			Start   struct {
				Line int `json:"line"`
			} `json:"start"`
			End struct {
				Line int `json:"line"`
			} `json:"end"`
			Extra struct {
				Message  string `json:"message"`
				Severity string `json:"severity"`
				Fix      string `json:"fix"`
			} `json:"extra"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	findings := make([]semgrepFinding, 0, len(result.Results))
	for _, r := range result.Results {
		findings = append(findings, semgrepFinding{
			CheckID:   r.CheckID,
			Path:      r.Path,
			StartLine: r.Start.Line,
			EndLine:   r.End.Line,
			Message:   r.Extra.Message,
			Severity:  r.Extra.Severity,
			Fix:       r.Extra.Fix,
		})
	}
	return findings, nil
}
