package appLogger

import (
	"fmt"

	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

func toZapFields(args ...string) []zap.Field {
	fields := make([]zap.Field, 0)

	for idx, arg := range args {
		argLabel := fmt.Sprintf("arg[%d]", idx)
		fields = append(fields, zap.String(argLabel, arg))
	}

	return fields
}

func Error(message string, args ...string) {
	logger.Error(message, toZapFields(args...)...)
}

func Info(message string, args ...string) {
	logger.Info(message, toZapFields(args...)...)
}

func Fatal(message string) {
	logger.Fatal("Fatal Error occured, app will halt", zap.String("Fatalpoint", message))
}
