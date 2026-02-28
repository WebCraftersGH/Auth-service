package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	LOGS_PATH     = "logs"
	LOGS_FILENAME = "log"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return nil
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		FullTimestamp: true,
	}

	_, err := os.Stat(LOGS_PATH)
	if os.IsNotExist(err) {
		os.MkdirAll(LOGS_PATH, 8644)
	} else {
		panic(err)
	}

	logFilename := LOGS_FILENAME + ".log"
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{logFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	e = logrus.NewEntry(l)
}
