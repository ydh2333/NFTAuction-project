package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger 初始化Zerolog
func InitLogger() {
	// 设置日志级别,默认info
	// level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	// if err != nil {
	// 	level = zerolog.InfoLevel
	// }
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// 美化日志输出
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
}
