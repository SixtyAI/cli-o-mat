package util

import (
	"fmt"
	"os"
)

func Fatal(err error) {
	fmt.Printf("%+v\n", err)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
