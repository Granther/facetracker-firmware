package errors

import (
	"fmt"
	"os"
)

// CheckError prints err to stderr and exits with code 1 if err is not nil. Otherwise it is no-op
func CheckError(err error) {
	if err != nil {	
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err) 
		os.Exit(1)
	}
}

// CheckErrorMsg prints err to stderr and exits with code 1 if err is not nil. errMsg is passed to the fprint output
func CheckErrorMsg(err error, errMsg string) {
	if err != nil {	
		fmt.Fprintf(os.Stderr, "An error occured while %s: %v\n", errMsg, err) 
		os.Exit(1)
	}
}

// Exit prints msg (with optional args), plus a newline, to stderr and then exits with code 1
func Exit(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

