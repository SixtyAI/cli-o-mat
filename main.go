package main

import (
	_ "embed"
	"fmt"

	"github.com/FasterBetter/cli-o-mat/cmd"
)

// nolint: stylecheck
//go:embed .version_string
var Version string

func main() {
	fmt.Printf("cli-o-mat v%s\n", Version) // nolint: forbidigo
	cmd.Execute()
}
