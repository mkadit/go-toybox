package logfile

import (
	"errors"
	"fmt"
	"io"
)

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{err, callers()}
}

// GetStack returns the stacktrace of the first error in err's tree has a stacktrace.
// The function returns true if it found such error, and false if there is no error with
// a stack in the err's tree.
func GetStack(err error) (StackTrace, bool) {
	var target *withStack
	ok := errors.As(err, &target)
	if !ok {
		return nil, false
	}
	return target.StackTrace(), true
}

// HasStack returns true if any error in err's tree has a stacktrace.
func HasStack(err error) bool {
	var target *withStack
	return errors.As(err, &target)
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
