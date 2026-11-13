package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"syscall"
	"time"
)

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

// runSemgrepConfig scans filePath with the rules at configPath and returns semgrep's JSON output.
func runSemgrepConfig(ctx context.Context, configPath, filePath string, extraFlags ...string) ([]byte, error) {
	args := append([]string{"scan", "--quiet", "--json"}, extraFlags...)
	args = append(args, "--config", configPath, filePath)

	cmd := exec.CommandContext(ctx, "semgrep", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = 2 * time.Second
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
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
