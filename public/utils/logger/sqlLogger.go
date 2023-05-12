package logger

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type SqlWriter struct {
	logger.Writer
}

func NewSqlWriter(w logger.Writer) *SqlWriter {
	return &SqlWriter{Writer: w}
}

func (w *SqlWriter) Printf(message string, data ...interface{}) {
	zap.L().Info(fmt.Sprintf(message+"\n", data...))
	w.Writer.Printf(message, data...)
}
