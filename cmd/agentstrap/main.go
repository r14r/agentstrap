package main

import (
	"os"

	"github.com/r14r/agentstrap/internal/cli"
)

// version is populated from VERSION by the build command.
var version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, version))
}
