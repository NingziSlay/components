package pkg

import (
	"fmt"
	"github.com/pkg/errors"
)

// ConfigError mapper 过程中会出现的错误
type ConfigError struct {
	msg string
	err error
}

// Error error interface
func (c ConfigError) Error() string {
	if c.err != nil {
		return errors.Wrap(c.err, c.msg).Error()
	}
	return c.msg
}

func newE(msg string, args ...interface{}) ConfigError {
	return ConfigError{msg: fmt.Sprintf(msg, args...)}
}

func wrapE(msg string, err error) ConfigError {
	return ConfigError{msg: msg, err: err}
}

// 预定义的一些异常
var (
	ErrorNonPointer = newE("mapStruct of non-pointer")
	ErrorNilInput   = newE("mapStruct of nil")
	ErrorNonStruct  = newE("mapStruct of non-struct ptr")
)
