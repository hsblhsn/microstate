package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/hsblhsn/microstate/state"
	"github.com/logrusorgru/aurora/v3"
)

type Logger struct {
	l io.Writer
	o io.Writer
}

func NewLogger() *Logger {
	return &Logger{
		l: os.Stderr,
		o: os.Stdout,
	}
}

func (l *Logger) Promotion(store *state.State, to state.ReleaseKind) {
	fromR := store.Latest(to)
	toR := store.Latest(to)
	fmt.Fprintf(l.l, "promoted %s to %s\n", aurora.BrightRed(fromR), aurora.BrightGreen(toR))
	fmt.Fprint(l.o, toR.String())
}

func (l *Logger) Error(v interface{}) {
	fmt.Fprintf(l.l, "ERROR: %v\n", aurora.BrightRed(v))
}

func (l *Logger) OK(v interface{}) {
	fmt.Fprintf(l.l, "OK: %v\n", aurora.BrightGreen(v))
}
