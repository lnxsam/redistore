package yerror

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

//	Error kinds based on grpc codes
const (
	KindNotFound        = codes.NotFound
	KindInvalidArgument = codes.InvalidArgument
	KindUnauthenticated = codes.Unauthenticated
	KindUnauthorized    = codes.PermissionDenied
	KindInternal        = codes.Internal
	KindUnexpected      = codes.Unknown

	LevelInfo  = logrus.InfoLevel
	LevelDebug = logrus.DebugLevel
	LevelWarn  = logrus.WarnLevel
	LevelError = logrus.ErrorLevel
)

// Op is the operation name that the error occurs at
type Op string

// Error is our project custom error type that implements go error
type Error struct {
	Op       Op
	Kind     codes.Code
	Err      error
	Severity logrus.Level
}

// E returns new Error filled with args
func E(args ...interface{}) *Error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case error:
			e.Err = arg
		case codes.Code:
			e.Kind = arg
		case logrus.Level:
			e.Severity = arg
		default:
			panic(fmt.Sprintf("bad call to E: %v", arg))
		}
	}
	return e
}

// Ops returns operations chain of an Error
func Ops(e *Error) []Op {
	res := []Op{e.Op}
	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}
	res = append(res, Ops(subErr)...)
	return res
}

// Kind returns kind of an Error
func Kind(err error) codes.Code {
	e, ok := err.(*Error)
	if !ok {
		return KindUnexpected
	}
	if e.Kind != 0 {
		return e.Kind
	}
	return Kind(e.Err)
}

// Level returns severity level of an Error
func Level(err error) logrus.Level {
	e, ok := err.(*Error)
	if !ok {
		return logrus.ErrorLevel
	}
	if e.Severity != 0 {
		return e.Severity
	}
	return Level(e.Err)
}

func (e *Error) Error() string {
	if subErr, ok := e.Err.(*Error); ok {
		return subErr.Error()
	}
	return e.Err.Error()
}
