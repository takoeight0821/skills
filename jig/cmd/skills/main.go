package main

import (
	"fmt"
	"os"

	"github.com/takoeight0821/skills/jig/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
