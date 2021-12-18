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
	from, err := to.Prev()
	if err != nil {
		fmt.Fprintf(l.l, "cli: could not determine previous release kind: %s\n", err)
		return
	}
	fromR := store.Latest(from)
	toR := store.Latest(to)
	fmt.Fprintf(
		l.l,
		"promoted %s(%s) to %s(%s)\n",
		aurora.BrightRed(fromR),
		fromR.BlockHash.Short(),
		aurora.BrightGreen(toR),
		toR.BlockHash.Short(),
	)
	fmt.Fprint(l.o, toR.String())
}

func (l *Logger) Error(v interface{}) {
	fmt.Fprintf(l.l, "ERROR: %v\n", aurora.BrightRed(v))
}

func (l *Logger) OK(v interface{}) {
	fmt.Fprintf(l.l, "OK: %v\n", aurora.BrightGreen(v))
}
