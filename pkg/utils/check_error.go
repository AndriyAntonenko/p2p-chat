package utils

import (
	"fmt"
	"os"
)

func HandleFatalError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Fatal error] %s. Error %s", msg, err.Error())
		os.Exit(1)
	}
}
