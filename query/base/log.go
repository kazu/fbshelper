package base

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/kazu/fbshelper/query/log"
)

const (
	FLT_NONE = 1 << iota
	FLT_NORMAL
	FLT_IS
	FLT_IO
)

type LogState struct {
	level  log.LogLevel
	filter byte
	writer io.Writer
}

var LogCurrentState LogState = LogState{level: log.LOG_ERROR, filter: FLT_NORMAL}

type L2Run interface{}

type LogOptParam struct {
	state  *LogState
	level  log.LogLevel
	filter byte
	logfns []func()
	fn     func() L2Run
}

type Log2Option func(*LogOptParam)

func L2OptFlag(level log.LogLevel, filter byte) Log2Option {
	return func(p *LogOptParam) {
		p.level = level
		p.filter = filter
	}
}

func L2OptF(fn LogFn) Log2Option {
	return func(p *LogOptParam) {
		p.logfns = append(p.logfns, func() { fn() })
	}
}

func L2OptRun(fn func() L2Run) Log2Option {

	return func(p *LogOptParam) {
		p.fn = fn
	}
}

func SetL2Current(level log.LogLevel, filter byte) {

	LogCurrentState.level = level
	LogCurrentState.filter = filter
	if LogCurrentState.level == LOG_DEBUG {
		os.Exit(0)
	}

}

func L2isEnable(param Log2Option) bool {
	logParam := LogOptParam{state: &LogCurrentState, filter: FLT_NONE, logfns: make([]func(), 0, 2)}

	param(&logParam)

	if logParam.state.level < logParam.level {
		return false
	}
	if logParam.state.filter&logParam.filter == 0 {
		return false
	}
	return true

}

// Log2 ... logger with lazy evaluation and filter mode
func Log2(opts ...Log2Option) (result interface{}) {
	logParam := LogOptParam{state: &LogCurrentState, filter: FLT_NONE, logfns: make([]func(), 0, 2)}

	for _, opt := range opts {
		opt(&logParam)
	}
	o := log.CurrentLogLevel

	enableAfterLog := false

	defer func() {
		if logParam.fn != nil {
			result = logParam.fn()
		}
		if enableAfterLog {
			SetLogLevel(logParam.state.level)
			neoLog(logParam.level, logParam.logfns[1])
			SetLogLevel(o)
		}
	}()

	if logParam.state.filter&logParam.filter == 0 {
		return
	}

	SetLogLevel(logParam.state.level)
	defer SetLogLevel(o)
	if len(logParam.logfns) > 0 {
		neoLog(logParam.level, logParam.logfns[0])
	}
	if len(logParam.logfns) >= 2 {
		enableAfterLog = true
	}
	return
}

var (
	L2_DEBUG_IS Log2Option = L2OptFlag(LOG_DEBUG, FLT_IS)
	L2_WARN_IS  Log2Option = L2OptFlag(LOG_WARN, FLT_IS)
	L2_WARN_IO  Log2Option = L2OptFlag(LOG_WARN, FLT_IO)
	L2_DEBUG_IO Log2Option = L2OptFlag(LOG_DEBUG, FLT_IO)
)

type LogFmtParam func() []interface{}

func L2fmt(s string, vlist ...interface{}) Log2Option {

	return L2OptF(func() LogArgs {
		FF(s, vlist...)
		return LogArgs{}
	})
}

// func L2Fmt(s string, fn LogFmtParam) Log2Option {

// 	return L2OptF(func() LogArgs {
// 		return F(s, fn())
// 	})
// }

func L2F(fn func()) Log2Option {

	return L2OptF(func() LogArgs {
		fn()
		return LogArgs{}
	})
}

func FF(s string, v ...interface{}) {
	fmt.Fprintf(LogCurrentState.writer, s, v...)
}

func neoLog(l LogLevel, fn func()) {

	if CurrentLogLevel < l {
		return
	}
	o := LogW

	var b strings.Builder
	LogW = &b

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
	fn()
	io.WriteString(LogW, b.String())
	LogW = o

}
