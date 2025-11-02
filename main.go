package main

import (
	"fmt"
	"os"

	"github.com/mesh-dell/github-activity/cmd"
)

func main() {
	if err := cmd.Execute(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
