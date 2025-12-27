package logging

import (
	"fmt"
	"os"
)

type StdoutLogger struct{}

func (StdoutLogger) Logf(format string, args ...any) {
	fmt.Fprintln(os.Stdout, format, args)

}

func (StdoutLogger) Errorf(err error, format string, args ...any) {
	if err != nil {
		format = format + ": %v"
		args = append(args, err.Error())
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
