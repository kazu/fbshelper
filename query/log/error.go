package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ERR_MUST_POINTER    error = errors.New("parameter must be pointer")
	ERR_INVALID_TYPE    error = errors.New("parameter invalid type")
	ERR_NOT_FOUND       error = errors.New("data is not found")
	ERR_READ_BUFFER     error = errors.New("cannot read least data")
	ERR_MORE_BUFFER     error = errors.New("require more data")
	ERR_NO_SUPPORT      error = errors.New("this method is not suppored")
	ERR_INVALID_INDEX   error = errors.New("invalid index number")
	ERR_NO_INCLUDE_ROOT error = errors.New("dosent include root buffer")
)

type LogLevel byte

var CurrentLogLevel LogLevel
var LogW io.Writer = os.Stderr

const (
	LOG_ERROR LogLevel = iota
	LOG_WARN
	LOG_DEBUG
)

type LogArgs struct {
	Fmt  string
	Infs []interface{}
}

type LogFn func() LogArgs //(string, interface{}...)

func SetLogLevel(l LogLevel) {
	CurrentLogLevel = l
}

func F(s string, v ...interface{}) LogArgs {
	return LogArgs{Fmt: s, Infs: v}
}

// if no output , not eval args
//  Log(LOG_DEBUG, func() LogArgs { return F("test %d \n", 1) })
func Log(l LogLevel, fn LogFn) {

	if CurrentLogLevel < l {
		return
	}

	var b strings.Builder
	switch l {
	case LOG_DEBUG:
		b.WriteString("D: ")
	case LOG_WARN:
		b.WriteString("W: ")
	case LOG_ERROR:
		b.WriteString("E: ")
	default:
		b.WriteString(" ")
	}
	args := fn()
	fmt.Fprintf(&b, args.Fmt, args.Infs...)
	io.WriteString(LogW, b.String())

	return

}
