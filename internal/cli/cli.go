// Package cli implements AgentStrap's command-line entry point.
package cli

import (
	"fmt"
	"io"
)

// Run executes the requested command and returns a process exit code.
// The CLI currently implements the bootstrap commands only. Management of
// clients, skills, agents, and MCP servers will be added incrementally.
func Run(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return invalidArguments(stderr)
		}
		printHelp(stdout)
		return 0
	case "version", "-v", "--version":
		if len(args) != 1 {
			return invalidArguments(stderr)
		}
		fmt.Fprintf(stdout, "agentstrap %s\n", version)
		return 0
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\nRun 'agentstrap --help' for usage.\n", args[0])
		return 2
	}
}

func invalidArguments(stderr io.Writer) int {
	fmt.Fprintln(stderr, "Unexpected arguments. Run 'agentstrap --help' for usage.")
	return 2
}

func printHelp(out io.Writer) {
	fmt.Fprint(out, "AgentStrap — One CLI. Every AI coding environment.\n\n")
	fmt.Fprint(out, "Usage:\n  agentstrap [command]\n\n")
	fmt.Fprint(out, "Commands:\n")
	fmt.Fprint(out, "  help       Display this help\n")
	fmt.Fprint(out, "  version    Show the version\n\n")
	fmt.Fprint(out, "Planned: client setup, skills, agents, MCP servers and config generation.\n")
}
