package util

import (
	"fmt"
	"os"
)

func Fatal(code int, err error) {
	fmt.Printf("%v\n", err)
	os.Exit(code)
}

func Fatalf(code int, format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(code)
}
