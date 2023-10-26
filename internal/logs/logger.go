package logs

import (
	"fmt"
	"io"
	"log"
	"os"
)

//nolint:gochecknoglobals
var (
	Debug Logger = discard{}
	Info  Logger = logger{log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)}
	Warn  Logger = logger{log.New(os.Stderr, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)}
)

type Logger interface {
	io.Writer
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

const callerSkip = 2

type logger struct{ l *log.Logger }

func NewDebug() Logger { //nolint:ireturn
	return logger{log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)}
}
func (l logger) Print(v ...interface{}) { _ = l.l.Output(callerSkip, fmt.Sprint(v...)) }
func (l logger) Printf(format string, v ...interface{}) {
	_ = l.l.Output(callerSkip, fmt.Sprintf(format, v...))
}

func (l logger) Write(p []byte) (n int, err error) {
	l.Print(string(p))
	return len(p), nil
}

type discard struct{}

func NewDiscard() Logger                            { return discard{} } //nolint:ireturn
func (l discard) Print(_ ...interface{})            {}
func (l discard) Printf(_ string, _ ...interface{}) {}
func (l discard) Write(p []byte) (n int, err error) { return len(p), nil }
