package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	debugLog = log.New(os.Stdin, "\033[30mDEBUG:\033[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stdin, "\033[31mERROR: \033[0m", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog  = log.New(os.Stdin, "\033[34mINFO: \033[0m", log.Ldate|log.Ltime|log.Lshortfile)
	loggers  = []*log.Logger{
		debugLog,
		errorLog,
		infoLog,
	}
	mu sync.Mutex
)

var (
	Debug  = debugLog.Println
	Debugf = debugLog.Printf
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func SetLogger(l *log.Logger) {
	mu.Lock()
	defer mu.Unlock()
	loggers = append(loggers, l)
}

const (
	DebugLevel = iota
	InfoLevel
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, l := range loggers {
		l.SetOutput(os.Stdout)
	}
	if level > ErrorLevel {
		errorLog.SetOutput(io.Discard)
	}
	if level > InfoLevel {
		infoLog.SetOutput(io.Discard)
	}
	if level > DebugLevel {
		debugLog.SetOutput(io.Discard)
	}

}
