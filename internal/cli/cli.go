// Package cli implements the AgentStrap command-line interface.
package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/r14r/agentstrap/internal/catalog"
)

type options struct {
	topic, client, path string
	apply               bool
}

func Run(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}
	var err error
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
	case "skill", "agent", "mcp":
		err = resourceCommand(args[0], args[1:], stdout)
	case "source":
		err = sourceCommand(args[1:], stdout)
	case "inspect":
		err = inspectCommand(args[1:], stdout)
	case "repo":
		if len(args) < 2 || args[1] != "inspect" {
			err = fmt.Errorf("usage: agentstrap repo inspect [--path DIR] [--client CLIENT]")
		} else {
			err = inspectCommand(args[2:], stdout)
		}
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\nRun 'agentstrap --help' for usage.\n", args[0])
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	return 0
}

func invalidArguments(stderr io.Writer) int {
	fmt.Fprintln(stderr, "Unexpected arguments. Run 'agentstrap --help' for usage.")
	return 2
}

func parseOptions(args []string) (options, []string, error) {
	opts := options{path: "."}
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--apply":
			opts.apply = true
		case arg == "--topic" || arg == "--client" || arg == "--path":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				return opts, nil, fmt.Errorf("%s requires a value", arg)
			}
			i++
			switch arg {
			case "--topic":
				opts.topic = args[i]
			case "--client":
				opts.client = strings.ToLower(args[i])
			case "--path":
				opts.path = args[i]
			}
		case strings.HasPrefix(arg, "--topic="):
			opts.topic = strings.TrimPrefix(arg, "--topic=")
		case strings.HasPrefix(arg, "--client="):
			opts.client = strings.ToLower(strings.TrimPrefix(arg, "--client="))
		case strings.HasPrefix(arg, "--path="):
			opts.path = strings.TrimPrefix(arg, "--path=")
		case strings.HasPrefix(arg, "-"):
			return opts, nil, fmt.Errorf("unknown option %q", arg)
		default:
			positional = append(positional, arg)
		}
	}
	return opts, positional, nil
}

func resourceCommand(kind string, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: agentstrap %s list|search|install", kind)
	}
	action := args[0]
	opts, positional, err := parseOptions(args[1:])
	if err != nil {
		return err
	}
	switch action {
	case "list":
		if len(positional) != 0 || opts.client != "" || opts.apply {
			return fmt.Errorf("usage: agentstrap %s list [--topic TOPIC]", kind)
		}
		return printItems(out, catalog.Search(kind, opts.topic, ""))
	case "search":
		if len(positional) == 0 || opts.client != "" || opts.apply {
			return fmt.Errorf("usage: agentstrap %s search QUERY [--topic TOPIC]", kind)
		}
		return printItems(out, catalog.Search(kind, opts.topic, strings.Join(positional, " ")))
	case "install":
		if len(positional) != 1 || opts.client == "" || opts.topic != "" {
			return fmt.Errorf("usage: agentstrap %s install ID --client codex|claude|gemini [--path DIR] [--apply]", kind)
		}
		if strings.Contains(positional[0], ":") {
			return installSourceItem(kind, positional[0], opts, out)
		}
		item, found := catalog.Find(kind, positional[0])
		if !found {
			return fmt.Errorf("unknown %s %q; use '%s list'", kind, positional[0], kind)
		}
		return installItem(item, opts, out)
	default:
		return fmt.Errorf("unknown %s action %q; use list, search or install", kind, action)
	}
}

func printItems(out io.Writer, items []catalog.Item) error {
	if len(items) == 0 {
		fmt.Fprintln(out, "No matches in the built-in catalog.")
		return nil
	}
	for _, item := range items {
		fmt.Fprintf(out, "%-6s %-20s [%s] %s\n", item.Kind, item.ID, strings.Join(item.Topics, ", "), item.Description)
	}
	return nil
}

func printHelp(out io.Writer) {
	fmt.Fprint(out, "AgentStrap — One CLI. Every AI coding environment.\n\n")
	fmt.Fprint(out, "Usage:\n  agentstrap [command]\n\nCommands:\n")
	fmt.Fprint(out, "  skill list|search|install   Find and install project skills\n")
	fmt.Fprint(out, "  agent list|search|install   Find and install agents\n")
	fmt.Fprint(out, "  mcp list|search|install     Find and configure MCP servers\n")
	fmt.Fprint(out, "  source list|show|add|remove|sync|search|install   Manage external GitHub registries\n")
	fmt.Fprint(out, "  inspect / repo inspect     Suggest project instruction and skill files\n")
	fmt.Fprint(out, "  version                    Show version\n")
	fmt.Fprint(out, "  help                       Display help\n\n")
	fmt.Fprint(out, "Flags: --topic TOPIC, --client codex|claude|gemini|all, --path DIR\n")
	fmt.Fprint(out, "Install previews changes by default; use --apply to write files.\n")
	fmt.Fprint(out, "Topics: development, video-editing, animation.\n")
}
