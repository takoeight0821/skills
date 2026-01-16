package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

type Logger struct {
	out io.Writer
	err io.Writer
}

var (
	infoColor    = color.New(color.FgGreen)
	warnColor    = color.New(color.FgYellow)
	errorColor   = color.New(color.FgRed)
	successColor = color.New(color.FgGreen, color.Bold)
)

func NewLogger(out, err io.Writer) *Logger {
	return &Logger{out: out, err: err}
}

func Default() *Logger {
	return &Logger{out: os.Stdout, err: os.Stderr}
}

func (l *Logger) Info(format string, args ...interface{}) {
	infoColor.Fprintf(l.err, "[INFO] ")
	fmt.Fprintf(l.err, format+"\n", args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	warnColor.Fprintf(l.err, "[WARN] ")
	fmt.Fprintf(l.err, format+"\n", args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	errorColor.Fprintf(l.err, "[ERROR] ")
	fmt.Fprintf(l.err, format+"\n", args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	successColor.Fprintf(l.err, format+"\n", args...)
}

func (l *Logger) Print(format string, args ...interface{}) {
	fmt.Fprintf(l.out, format+"\n", args...)
}
