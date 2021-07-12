package hook

import (
	"context"
	"time"

	"github.com/nori-io/common/v5/pkg/domain/logger"
	gormLogger "gorm.io/gorm/logger"
)

type Logger struct {
	Origin logger.FieldLogger
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l.LogMode(level)
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	l.Info(ctx, s)
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.Warn(ctx, s)
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	l.Error(ctx, s)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	l.Trace(ctx, begin, fc, err)
}
