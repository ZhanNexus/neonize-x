// Copyright (c) 2021 Tulir Asokan
//

// Package waLog contains a simple logger interface used by the other whatsmeow packages.
package utils


/*

   #cgo CFLAGS: -I${SRCDIR}/../header -I${SRCDIR}/../python
   #include <stdlib.h>
   #include <stdbool.h>
   #include <stdint.h>
   #include <string.h>
   #include "cstruct.h"
   #include "pythonptr.h"
*/
import (
	"C"
	defproto "github.com/krypton-byte/neonize/defproto"
	"google.golang.org/protobuf/proto"
	"unsafe"
)
import (
	"fmt"
	"strings"
	"time"
)


func getBytesAndSize(data []byte) (*C.char, C.size_t) {
	messageSourceCDATA := (*C.char)(unsafe.Pointer(&data[0]))
	messageSourceCSize := C.size_t(len(data))
	return messageSourceCDATA, messageSourceCSize
}

// Logger is a simple logger interface that can have subloggers for specific areas.
type Logger interface {
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Sub(module string) Logger
}

type noopLogger struct{}

func (n *noopLogger) Errorf(_ string, _ ...interface{}) {}
func (n *noopLogger) Warnf(_ string, _ ...interface{})  {}
func (n *noopLogger) Infof(_ string, _ ...interface{})  {}
func (n *noopLogger) Debugf(_ string, _ ...interface{}) {}
func (n *noopLogger) Sub(_ string) Logger               { return n }

// Noop is a no-op Logger implementation that silently drops everything.
var Noop Logger = &noopLogger{}

type stdoutLogger struct {
	mod   string
	color bool
	min   int
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
		message: msg,
		level: level,
		name: s.mod,
	}
	buff, err := proto.Marshal(&log_msg)
	if err != nil {
		panic(err)
	}
	uchars, size := getBytesAndSize(buff)
	C.call_c_func_callback_bytes(callback, uchars, size)
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
func (s *stdoutLogger) Sub(mod string) Logger {
	return &stdoutLogger{mod: fmt.Sprintf("%s/%s", s.mod, mod), color: s.color, min: s.min}
}

// Stdout is a simple Logger implementation that outputs to stdout. The module name given is included in log lines.
//
// minLevel specifies the minimum log level to output. An empty string will output all logs.
//
// If color is true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes.
func NewLogger(module string, minLevel string, callback C.ptr_to_python_function_bytes) Logger {
	return &stdoutLogger{mod: module, min: levelToInt[strings.ToUpper(minLevel)], callback: callback}
}