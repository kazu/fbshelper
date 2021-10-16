package base

import (
	log "github.com/kazu/fbshelper/query/log"
)

const (
	FLT_NONE = 1 << iota
	FLT_NORMAL
	FLT_IS
)

type LogState struct {
	level  log.LogLevel
	filter byte
}

var LogCurrentState LogState = LogState{level: log.LOG_ERROR, filter: FLT_NORMAL}

type L2Run interface{}

type LogOptParam struct {
	state  *LogState
	level  log.LogLevel
	filter byte
	logfns []LogFn
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
		p.logfns = append(p.logfns, fn)
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

}

// Log2 ... logger with lazy evaluation and filter mode
func Log2(opts ...Log2Option) (result interface{}) {
	logParam := LogOptParam{state: &LogCurrentState, filter: FLT_NONE, logfns: make([]LogFn, 0, 2)}

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
			Log(logParam.level, logParam.logfns[1])
			SetLogLevel(o)
		}
	}()

	if logParam.state.filter&logParam.filter == 0 {
		return
	}

	SetLogLevel(logParam.state.level)
	defer SetLogLevel(o)
	if len(logParam.logfns) > 0 {
		Log(logParam.level, logParam.logfns[0])
	}
	if len(logParam.logfns) >= 2 {
		enableAfterLog = true
	}
	return
}

var (
	L2_DEBUG_IS Log2Option = L2OptFlag(LOG_DEBUG, FLT_IS)
)

func L2fmt(s string, vlist ...interface{}) Log2Option {

	return L2OptF(func() LogArgs {
		return F(s, vlist...)
	})
}
