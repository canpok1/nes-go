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
	LogLevel Level
}

var config = Config{
	LogLevel: LevelInfo,
}

// SetOutput ...
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// print ...
func print(level string, format string, v ...interface{}) {
	f := "[%v]" + format
	args := append([]interface{}{level}, v...)
	log.Printf(f, args...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	if config.LogLevel > LevelFatal {
		return
	}
	print("F", format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	if config.LogLevel > LevelWarn {
		return
	}
	print("W", format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
	if config.LogLevel > LevelInfo {
		return
	}
	print("I", format, v...)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	if config.LogLevel > LevelDebug {
		return
	}
	print("D", format, v...)
}
