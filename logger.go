package pkg

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var l zerolog.Logger

// InitLogger 创建 l 实例
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
	l = zerolog.New(output).With().Caller().Timestamp().Logger().Level(level)
}

// GetLogger 获取 l 实例，如果 l 为空，则初始化一个 l 实例
func GetLogger() zerolog.Logger {
	if &l != nil {
		return l
	}
	InitLogger(false)
	return l
}

// WithLoggerLevel 返回一个自定义日志等级的 l 实例
func WithLoggerLevel(level zerolog.Level) zerolog.Logger {
	if &l != nil {
		return l.Level(level)
	}
	InitLogger(false)
	return l.Level(level)
}

// WithSampleLog 返回一个 sampling log 实例
// debug 等级每秒输出前 5 条日志，之后每 20 条输出一条
// info 等级每秒输出前 5 条日志，之后每 10 条输出一条
func WithSampleLog() zerolog.Logger {
	if &l == nil {
		InitLogger(false)
	}
	return l.Sample(zerolog.LevelSampler{
		DebugSampler: &zerolog.BurstSampler{
			Burst:       5,
			Period:      time.Second * 1,
			NextSampler: &zerolog.BasicSampler{N: 20},
		},
		InfoSampler: &zerolog.BurstSampler{
			Burst:       5,
			Period:      time.Second * 1,
			NextSampler: &zerolog.BasicSampler{N: 10},
		},
	})
}
