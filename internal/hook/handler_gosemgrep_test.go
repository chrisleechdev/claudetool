package hook

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/laurenkt/claudetool/internal/hook/semgreprules"
)

func requireSemgrep(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("semgrep"); err != nil {
		t.Skip("semgrep not installed")
	}
}

func writeGoFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestGoSemgrepRules(t *testing.T) {
	requireSemgrep(t)

	rulesDir, err := filepath.Abs("semgreprules")
	if err != nil {
		t.Fatal(err)
	}
	fixturesDir := filepath.Join(rulesDir, "testdata")

	entries, err := fs.ReadDir(semgreprules.FS, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		ruleID := strings.TrimSuffix(entry.Name(), ".yml")
		t.Run(ruleID, func(t *testing.T) {
			if _, err := os.Stat(filepath.Join(fixturesDir, ruleID+".go")); err != nil {
				t.Errorf("rule %s has no fixture: %v", ruleID, err)
			}
		})
	}

	// One semgrep run for all rules; each spawn costs seconds
	out, err := exec.Command("semgrep", "--test", "--metrics=off", "--disable-version-check", "--config", rulesDir, fixturesDir).CombinedOutput()
	if err != nil {
		t.Fatalf("semgrep --test failed: %v\n%s", err, out)
	}
}

func TestGoSemgrepBlocksWrongErrCheck(t *testing.T) {
	requireSemgrep(t)

	path := writeGoFile(t, "bad.go", `package sample

import "fmt"

func handle() error {
	err := fmt.Errorf("a")
	if someErr := something(); err != nil {
		return someErr
	}
	return nil
}

func something() error { return nil }
`)

	input := makeToolInput("PostToolUse", "Write", WriteInput{
		FilePath: path,
		Content:  "ignored",
	})
	out := runHandlerOutput(t, "go-semgrep", input)
	if out == nil || out.Decision != "block" {
		t.Fatalf("output = %+v, want block", out)
	}
	if !strings.Contains(out.Reason, "[wrong-err-check]") {
		t.Errorf("reason = %q, want wrong-err-check finding", out.Reason)
	}
}

func TestGoSemgrepInfoFindingIsAdvisory(t *testing.T) {
	requireSemgrep(t)

	path := writeGoFile(t, "clocked.go", `package sample

import "time"

func now() time.Time {
	return time.Now()
}
`)

	input := makeToolInput("PostToolUse", "Write", WriteInput{
		FilePath: path,
		Content:  "ignored",
	})
	out := runHandlerOutput(t, "go-semgrep", input)
	if out == nil {
		t.Fatal("output = nil, want advisory context")
	}
	if out.Decision != "" {
		t.Errorf("decision = %q, want no block for INFO findings", out.Decision)
	}
	if out.HookSpecificOutput == nil || !strings.Contains(out.HookSpecificOutput.AdditionalContext, "[discourage-time-now]") {
		t.Errorf("output = %+v, want discourage-time-now in additional context", out)
	}
}

func TestGoSemgrepEditOnlyReportsChangedLines(t *testing.T) {
	requireSemgrep(t)

	path := writeGoFile(t, "legacy.go", `package sample

import (
	"context"

	"github.com/monzo/slog"
)

func legacy(ctx context.Context) {
	slog.WithParam(ctx, "key", "value")
}

func added(ctx context.Context) context.Context {
	return ctx
}
`)

	cleanEdit := makeToolInput("PostToolUse", "Edit", EditInput{
		FilePath:  path,
		OldString: "return nil",
		NewString: "func added(ctx context.Context) context.Context {\n\treturn ctx\n}",
	})
	if out := runHandlerOutput(t, "go-semgrep", cleanEdit); out != nil {
		t.Errorf("output = %+v, want nil: the legacy finding is outside the edited lines", out)
	}

	badEdit := makeToolInput("PostToolUse", "Edit", EditInput{
		FilePath:  path,
		OldString: "slog.WithParams(ctx, nil)",
		NewString: `slog.WithParam(ctx, "key", "value")`,
	})
	out := runHandlerOutput(t, "go-semgrep", badEdit)
	if out == nil || out.Decision != "block" {
		t.Fatalf("output = %+v, want block for the edited line", out)
	}
	if !strings.Contains(out.Reason, "[slog-withparam-return-not-used]") {
		t.Errorf("reason = %q, want slog-withparam-return-not-used finding", out.Reason)
	}
}

func TestGoSemgrepTestOnlyRuleSkipsNonTestFiles(t *testing.T) {
	requireSemgrep(t)

	content := `package sample

import (
	"testing"

	"github.com/monzo/terrors"
)

func check(t *testing.T, err error) {
	terrors.Is(err, terrors.ErrInternalService, "oops")
}
`

	nonTestInput := makeToolInput("PostToolUse", "Write", WriteInput{
		FilePath: writeGoFile(t, "check.go", content),
		Content:  "ignored",
	})
	if out := runHandlerOutput(t, "go-semgrep", nonTestInput); out != nil {
		t.Errorf("output = %+v, want nil: terrors-is-not-checked only applies to _test.go", out)
	}

	testInput := makeToolInput("PostToolUse", "Write", WriteInput{
		FilePath: writeGoFile(t, "check_test.go", content),
		Content:  "ignored",
	})
	out := runHandlerOutput(t, "go-semgrep", testInput)
	if out == nil || out.Decision != "block" {
		t.Fatalf("output = %+v, want block", out)
	}
	if !strings.Contains(out.Reason, "[terrors-is-not-checked]") {
		t.Errorf("reason = %q, want terrors-is-not-checked finding", out.Reason)
	}
}

func TestGoSemgrepSkipsWithoutSpawning(t *testing.T) {
	binDir := t.TempDir()
	marker := filepath.Join(t.TempDir(), "spawned")
	fakeSemgrep := "#!/bin/sh\ntouch \"$GO_SEMGREP_MARKER\"\necho '{\"results\":[]}'\n"
	if err := os.WriteFile(filepath.Join(binDir, "semgrep"), []byte(fakeSemgrep), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GO_SEMGREP_MARKER", marker)

	skipped := []string{
		"/src/README.md",
		"/src/service.py",
		"/src/service.pb.go",
		"/src/vendor/github.com/foo/bar.go",
		"vendor/github.com/foo/bar.go",
	}
	for _, path := range skipped {
		t.Run(path, func(t *testing.T) {
			input := makeToolInput("PostToolUse", "Write", WriteInput{
				FilePath: path,
				Content:  "ignored",
			})
			if out := runHandlerOutput(t, "go-semgrep", input); out != nil {
				t.Errorf("output = %+v, want nil", out)
			}
			if _, err := os.Stat(marker); err == nil {
				t.Errorf("semgrep was spawned for %s", path)
			}
		})
	}

	input := makeToolInput("PostToolUse", "Write", WriteInput{
		FilePath: "/src/service.go",
		Content:  "ignored",
	})
	runHandlerOutput(t, "go-semgrep", input)
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("semgrep was not spawned for a .go file: %v", err)
	}
}

func TestChangedLineRanges(t *testing.T) {
	path := writeGoFile(t, "changed.go", "package sample\n\nfunc a() {}\n\nfunc b() {\n\tx := 1\n}\n\nfunc a() {}\n")

	tests := []struct {
		name  string
		tool  string
		input any
		want  []lineRange
	}{
		{
			name:  "write is the whole file",
			tool:  "Write",
			input: WriteInput{FilePath: path},
			want:  nil,
		},
		{
			name:  "edit covers every occurrence",
			tool:  "Edit",
			input: EditInput{FilePath: path, NewString: "func a() {}", ReplaceAll: true},
			want:  []lineRange{{start: 3, end: 3}, {start: 9, end: 9}},
		},
		{
			name:  "multiline edit spans its lines",
			tool:  "Edit",
			input: EditInput{FilePath: path, NewString: "func b() {\n\tx := 1\n}\n"},
			want:  []lineRange{{start: 5, end: 7}},
		},
		{
			name:  "edit whose new text is missing is the whole file",
			tool:  "Edit",
			input: EditInput{FilePath: path, NewString: "func missing() {}"},
			want:  nil,
		},
		{
			name:  "deletion is the whole file",
			tool:  "Edit",
			input: EditInput{FilePath: path, NewString: ""},
			want:  nil,
		},
		{
			name: "multi-edit is the union of its edits",
			tool: "MultiEdit",
			input: MultiEditInput{
				FilePath: path,
				Edits: []EditInput{
					{NewString: "func b() {"},
					{NewString: "x := 1"},
				},
			},
			want: []lineRange{{start: 5, end: 5}, {start: 6, end: 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolInput, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			in := &Input{ToolName: tt.tool, ToolInput: toolInput}
			if got := changedLineRanges(in, path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ranges = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnChangedLines(t *testing.T) {
	findings := []semgrepFinding{
		{CheckID: "before", StartLine: 1, EndLine: 2},
		{CheckID: "straddles", StartLine: 4, EndLine: 6},
		{CheckID: "inside", StartLine: 7, EndLine: 7},
		{CheckID: "after", StartLine: 20, EndLine: 21},
	}
	ranges := []lineRange{{start: 5, end: 8}}

	var got []string
	for _, finding := range onChangedLines(findings, ranges) {
		got = append(got, finding.CheckID)
	}
	want := []string{"straddles", "inside"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("kept = %v, want %v", got, want)
	}

	if all := onChangedLines(findings, nil); len(all) != len(findings) {
		t.Errorf("no ranges kept %d findings, want all %d", len(all), len(findings))
	}
}
