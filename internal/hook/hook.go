package hook

import (
	"github.com/jinzhu/gorm"
	"github.com/nori-io/common/v5/pkg/domain/logger"
)

type Logger struct {
	Origin logger.FieldLogger
}

func (l *Logger) Print(values ...interface{}) {
	for _, v := range gorm.LogFormatter(values...) {
		l.Origin.Debug("%s\n", v)
	}
}
