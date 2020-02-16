package log

import (
	"io"
	"log"
	"runtime"
)

// Level ...
type Level int

const (
	// LevelTrace ...
	LevelTrace Level = iota
	// LevelDebug ...
	LevelDebug
	// LevelInfo ...
	LevelInfo
	// LevelWarn ...
	LevelWarn
	// LevelFatal ...
	LevelFatal
)

// Config ...
type Config struct {
	logLevel         Level
	enableLevelLabel bool
	enableFuncName   bool
}

var config = Config{
	logLevel:         LevelDebug,
	enableLevelLabel: true,
	enableFuncName:   true,
}

// SetOutput ...
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// SetLogLevel ...
func SetLogLevel(l Level) {
	config.logLevel = l
}

// SetEnableLevelLabel ...
func SetEnableLevelLabel(enable bool) {
	config.enableLevelLabel = enable
}

// SetEnableFuncName ...
func SetEnableFuncName(enable bool) {
	config.enableFuncName = enable
}

// SetEnableTimestamp ...
func SetEnableTimestamp(enable bool) {
	if enable {
		log.SetFlags(log.Flags() | log.Ldate | log.Ltime)
	} else {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}
}

// print ...
func print(level string, format string, v ...interface{}) {
	args := v
	f := format

	if config.enableFuncName {
		pc, _, _, ok := runtime.Caller(2)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()

			args = append([]interface{}{funcName}, args...)
			f = "[%s]" + f
		}
	}

	if config.enableLevelLabel {
		args = append([]interface{}{level}, args...)
		f = "[%v]" + f
	}

	log.Printf(f, args...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	if config.logLevel > LevelFatal {
		return
	}
	print("F", format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	if config.logLevel > LevelWarn {
		return
	}
	print("W", format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
	if config.logLevel > LevelInfo {
		return
	}
	print("I", format, v...)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	if config.logLevel > LevelDebug {
		return
	}
	print("D", format, v...)
}

// Trace ...
func Trace(format string, v ...interface{}) {
	if config.logLevel > LevelTrace {
		return
	}
	print("T", format, v...)
}
