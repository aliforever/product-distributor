package distributor

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Error(err error)
	Errorf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

type DefaultLogger struct {
	output io.Writer
}

func NewDefaultLogger(output io.Writer) *DefaultLogger {
	if output == nil {
		output = os.Stdout
	}

	return &DefaultLogger{
		output: output,
	}
}

func (dl *DefaultLogger) Error(err error) {
	err = fmt.Errorf("error: %s", err)

	_, _ = fmt.Fprintln(dl.output, err.Error())
}

func (dl *DefaultLogger) Errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("error: %s", format)

	_, _ = fmt.Fprintf(dl.output, format, args...)
}

func (dl *DefaultLogger) Info(args ...interface{}) {
	args = append([]interface{}{"info:"}, args...)

	_, _ = fmt.Fprintln(dl.output, args...)
}

func (dl *DefaultLogger) Infof(format string, args ...interface{}) {
	format = fmt.Sprintf("info: %s", format)

	_, _ = fmt.Fprintf(dl.output, format, args...)
}
