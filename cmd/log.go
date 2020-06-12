package cmd

import "fmt"

var (
	verbose bool
)

func logf(f string, _v ...interface{}) {
	if verbose {
		fmt.Printf(f+"\n", _v...)
	}
}

func log(_v ...interface{}) {
	if verbose {
		fmt.Println(_v...)
	}
}
