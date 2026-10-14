package main

import (
	"os"

	"github.com/tiagozaccaro/ai-project-starter/internal/cli"
)

// osExit is the function called when the CLI exits. Defaults to os.Exit
// but can be overridden in tests.
var osExit = os.Exit

// Run is the main entry point that can be tested.
func Run() {
	cli.Execute()
}

func main() {
	Run()
	osExit(0)
}
