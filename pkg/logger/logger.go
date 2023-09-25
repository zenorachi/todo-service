package logger

import (
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func init() {
	slog.SetDefault(logger)
}

func Info(service string, args ...interface{}) {
	slog.Info("unit "+service, "status", args)
}

func Error(service string, args ...interface{}) {
	slog.Error("service: "+service, "status", args)
}

func Fatal(service string, args ...interface{}) {
	slog.Error("service: "+service, "status", args)
	os.Exit(1)
}
