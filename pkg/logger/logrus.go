package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

func Init(out io.Writer, env string) {
	if env == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	logrus.SetOutput(out)
}

func Debug(msg ...interface{}) {
	// file, line string,
	// logrus.WithFields(logrus.Fields{
	// 	"file": file,
	// 	"line": line,
	// }).Debug(msg...)
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Error(msg ...interface{}) {
	// pkg, file, function, err string,
	// logrus.WithFields(logrus.Fields{
	// 	"package":  pkg,
	// 	"file":     file,
	// 	"function": function,
	// 	"error":    err,
	// }).Error(msg...)
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	logrus.Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
