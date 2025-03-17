package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log – глобальный логгер, который можно использовать во всём приложении.
var Log = logrus.New()

// Init инициализирует настройки логгера.
func Init() {
	Log.Out = os.Stdout
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
