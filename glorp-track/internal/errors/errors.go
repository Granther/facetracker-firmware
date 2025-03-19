package errors

import (
	"os"
	"fmt"
)

func ProcessError(code int, msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Code: %v, %s: %v\n", code, msg, err))
		os.Exit(1)
	}
}
