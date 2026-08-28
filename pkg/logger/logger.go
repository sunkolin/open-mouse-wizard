package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	warnLogger  = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)

func Info(format string, args ...interface{}) {
	infoLogger.Println(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	warnLogger.Println(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	errorLogger.Println(fmt.Sprintf(format, args...))
}
