package main

import (
	"github.com/bu3/rebrickable-cli/cli/cmd"
	"os"
)

func main() {
	exitCode := cmd.Execute()
	os.Exit(exitCode)
}
