package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		exitCode   int
		wantOutput string
		wantError  string
	}{
		{name: "default help", wantOutput: "Usage:", exitCode: 0},
		{name: "explicit help", args: []string{"help"}, wantOutput: "Commands:", exitCode: 0},
		{name: "help flag", args: []string{"--help"}, wantOutput: "agentstrap [command]", exitCode: 0},
		{name: "version", args: []string{"version"}, wantOutput: "agentstrap 1.2.3", exitCode: 0},
		{name: "version flag", args: []string{"--version"}, wantOutput: "agentstrap 1.2.3", exitCode: 0},
		{name: "unknown command", args: []string{"mcp"}, wantError: "Unknown command: mcp", exitCode: 2},
		{name: "unexpected arguments", args: []string{"version", "extra"}, wantError: "Unexpected arguments", exitCode: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			got := Run(tt.args, &stdout, &stderr, "1.2.3")
			if got != tt.exitCode {
				t.Fatalf("exit code = %d, want %d", got, tt.exitCode)
			}
			if !strings.Contains(stdout.String(), tt.wantOutput) {
				t.Errorf("stdout = %q, expected %q", stdout.String(), tt.wantOutput)
			}
			if !strings.Contains(stderr.String(), tt.wantError) {
				t.Errorf("stderr = %q, expected %q", stderr.String(), tt.wantError)
			}
		})
	}
}
