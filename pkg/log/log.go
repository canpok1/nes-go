package log

import (
	"io"
	"log"
)

// Level ...
type Level int

const (
	// LevelDebug ...
	LevelDebug Level = iota
	// LevelInfo ...
	LevelInfo
	// LevelWarn ...
	LevelWarn
	// LevelFatal ...
	LevelFatal
)

// Config ...
type Config struct {
	logLevel Level
}

var config = Config{
	logLevel: LevelDebug,
}

// SetOutput ...
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// SetLogLevel ...
func SetLogLevel(l Level) {
	config.logLevel = l
}

// print ...
func print(level string, format string, v ...interface{}) {
	f := "[%v]" + format
	args := append([]interface{}{level}, v...)
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
