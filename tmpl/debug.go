// +build debug

package tmpl

import (
	"fmt"
	"os"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.Create("tmpl.log")
	if err != nil {
		panic(err)
	}
}

func debugf(f string, args ...interface{}) {
	fmt.Fprintf(logFile, f, args...)
}
