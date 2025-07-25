// Copyright (c) 2021 Tulir Asokan
//

// Package waLog contains a simple logger interface used by the other whatsmeow packages.
package utils


/*

   #include <stdlib.h>
   #include <stdbool.h>
   #include <stdint.h>
   #include <string.h>
   #include "../header/cstruct.h"
   #include "../python/pythonptr.h"
*/
import "C"

import (
	defproto "github.com/krypton-byte/neonize/defproto"
	"google.golang.org/protobuf/proto"
	waLog "go.mau.fi/whatsmeow/util/log"
	"unsafe"
)
import (
	"fmt"
	"strings"
	// "time"
)


func getBytesAndSize(data []byte) (*C.char, C.size_t) {
	messageSourceCDATA := (*C.char)(unsafe.Pointer(&data[0]))
	messageSourceCSize := C.size_t(len(data))
	return messageSourceCDATA, messageSourceCSize
}

// Logger is a simple logger interface that can have subloggers for specific areas.
type Logger interface {
	waLog.Logger
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Errorf(_ string, _ ...interface{}) {}
func (n *noopLogger) Warnf(_ string, _ ...interface{})  {}
func (n *noopLogger) Infof(_ string, _ ...interface{})  {}
func (n *noopLogger) Debugf(_ string, _ ...interface{}) {}
func (n *noopLogger) Sub(_ string) waLog.Logger               { return n }

// Noop is a no-op Logger implementation that silently drops everything.
var Noop Logger = &noopLogger{}

type stdoutLogger struct {
	mod   string
	min   int
	
	callback  C.ptr_to_python_function_callback_bytes2
}

var colors = map[string]string{
	"INFO":  "\033[36m",
	"WARN":  "\033[33m",
	"ERROR": "\033[31m",
}

var levelToInt = map[string]int{
	"":      -1,
	"DEBUG": 0,
	"INFO":  1,
	"WARN":  2,
	"ERROR": 3,
}

func (s *stdoutLogger) outputf(level, msg string, args ...interface{}) {
	if levelToInt[level] < s.min {
		return
	}
	log_msg := defproto.LogEntry{
		Message: proto.String(msg),
		Level:   proto.String(level),
		Name:    proto.String(s.mod),
	}
	buff, err := proto.Marshal(&log_msg)
	if err != nil {
		panic(err)
	}
	uchars, size := getBytesAndSize(buff)
	C.call_c_func_callback_bytes2(s.callback, uchars, size)
	// var colorStart, colorReset string
	// if s.color {
		//  colorStart = colors[level]
		// colorReset = "\033[0m"
	// }
	// fmt.Printf("%s%s [%s %s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, s.mod, level, fmt.Sprintf(msg, args...), colorReset)
}

func (s *stdoutLogger) Errorf(msg string, args ...interface{}) { s.outputf("ERROR", msg, args...) }
func (s *stdoutLogger) Warnf(msg string, args ...interface{})  { s.outputf("WARN", msg, args...) }
func (s *stdoutLogger) Infof(msg string, args ...interface{})  { s.outputf("INFO", msg, args...) }
func (s *stdoutLogger) Debugf(msg string, args ...interface{}) { s.outputf("DEBUG", msg, args...) }
func (s *stdoutLogger) Sub(mod string) waLog.Logger {
	return &stdoutLogger{mod: fmt.Sprintf("%s/%s", s.mod, mod), callback: s.callback, min: s.min}
}

// Stdout is a simple Logger implementation that outputs to stdout. The module name given is included in log lines.
//
// minLevel specifies the minimum log level to output. An empty string will output all logs.
//
// If color is true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes.
func NewLogger(module string, minLevel string, callback C.ptr_to_python_function_callback_bytes2) Logger {
	return &stdoutLogger{mod: module, min: levelToInt[strings.ToUpper(minLevel)], callback: callback}
}