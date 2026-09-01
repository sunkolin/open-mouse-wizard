package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

type safeWriter struct {
	w io.Writer
}

func (s safeWriter) Write(p []byte) (int, error) {
	n, _ := s.w.Write(p)
	if n == 0 {
		return len(p), nil
	}
	return n, nil
}

func init() {
	logDir := getLogDir()
	os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "mouse-wizard.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		f, _ = os.OpenFile("mouse-wizard.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
	logFile = f

	var writers []io.Writer
	if f != nil {
		writers = append(writers, safeWriter{os.Stdout})
		writers = append(writers, f)
	} else {
		writers = append(writers, safeWriter{os.Stdout})
	}
	mw := io.MultiWriter(writers...)

	infoLogger = log.New(mw, "[INFO] ", log.LstdFlags|log.Lmicroseconds)
	warnLogger = log.New(mw, "[WARN] ", log.LstdFlags|log.Lmicroseconds)

	var errWriters []io.Writer
	if f != nil {
		errWriters = append(errWriters, safeWriter{os.Stderr})
		errWriters = append(errWriters, f)
	} else {
		errWriters = append(errWriters, safeWriter{os.Stderr})
	}
	errMw := io.MultiWriter(errWriters...)
	errorLogger = log.New(errMw, "[ERROR] ", log.LstdFlags|log.Lmicroseconds)
}

func getLogDir() string {
	exePath, err := os.Executable()
	if err == nil {
		return filepath.Dir(exePath)
	}
	wd, _ := os.Getwd()
	return wd
}

func LogPath() string {
	if logFile != nil {
		return logFile.Name()
	}
	return ""
}

func Info(format string, args ...interface{}) {
	infoLogger.Println(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	warnLogger.Println(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	errorLogger.Println(fmt.Sprintf(format, args...))
}

func SystemInfo() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Runtime: %s/%s (%s)", runtime.GOOS, runtime.GOARCH, runtime.Version()))
	sb.WriteString(fmt.Sprintf(" PID=%d", os.Getpid()))

	if runtime.GOOS == "windows" {
		if vi := windows.RtlGetVersion(); vi != nil {
			product := productName(uint32(vi.ProductType))
			sb.WriteString(fmt.Sprintf(" OS=Windows %d.%d.%d %s", vi.MajorVersion, vi.MinorVersion, vi.BuildNumber, product))
		}
	}

	if hostname, err := os.Hostname(); err == nil {
		sb.WriteString(fmt.Sprintf(" Host=%s", hostname))
	}

	if wd, err := os.Getwd(); err == nil {
		sb.WriteString(fmt.Sprintf(" CWD=%s", wd))
	}

	if user := os.Getenv("USERNAME"); user != "" {
		sb.WriteString(fmt.Sprintf(" User=%s", user))
	}

	return sb.String()
}

func productName(productType uint32) string {
	switch productType {
	case 0:
		return "Professional"
	case 1:
		return "Datacenter"
	case 2:
		return "Enterprise"
	case 3:
		return "Web"
	default:
		return fmt.Sprintf("ProductType=%d", productType)
	}
}

func WriteDump() {
	Info("=== System Dump Start ===")
	Info("  %s", SystemInfo())
	if wd, err := os.Getwd(); err == nil {
		Info("  Working Dir: %s", wd)
	}
	if exe, err := os.Executable(); err == nil {
		Info("  Executable: %s", exe)
	}
	Info("  CPU Count: %d", runtime.NumCPU())
	if logFile != nil {
		Info("  Log File: %s", logFile.Name())
	}
	Info("=== System Dump End ===")
}
