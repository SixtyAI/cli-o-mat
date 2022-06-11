package main

import (
	_ "embed"
	"fmt"
)

// nolint: stylecheck
//go:embed .version_string
var Version string

func main() {
	fmt.Printf("cli-o-mat v%s\n", Version) // nolint: forbidigo
}
