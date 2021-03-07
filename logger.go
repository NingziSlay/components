package components

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// InitLogger 创建 log 实例
func InitLogger(debug bool) {
	var (
		level            = zerolog.InfoLevel
		output io.Writer = os.Stderr
	)
	if debug {
		level = zerolog.DebugLevel
		output = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.FormatLevel = func(i interface{}) string {
				return fmt.Sprintf("[%s]", i.(string))
			}
			w.FormatLevel = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("|  %-6s|", i))
			}
			w.FormatFieldValue = func(i interface{}) string {
				return fmt.Sprintf("\"%s\"", i)
			}
		})
	}
	zerolog.TimeFieldFormat = time.RFC3339
	//zerolog.SetGlobalLevel(level)
	log = zerolog.New(output).With().Caller().Timestamp().Logger().Level(level)
	fmt.Println(log.GetLevel())
}

// GetLogger 获取 log 实例，如果 log 为空，则初始化一个 log 实例
func GetLogger() zerolog.Logger {
	if &log != nil {
		return log
	}
	InitLogger(false)
	return log
}

// SetLoggerLevel 返回一个自定义日志等级的 log 实例
func WithLoggerLevel(level zerolog.Level) zerolog.Logger {
	if &log != nil {
		return log.Level(level)
	}
	InitLogger(false)
	return log.Level(level)
}
