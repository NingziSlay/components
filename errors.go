package components

import (
	"errors"
	"fmt"
)

type ConfigError struct {
}

var (
	ErrorNotWriteable    = errors.New("dest is not writable, try use ptr")
	ErrorArrayOutOfRange = errors.New("array out of range")
	ErrorUnsupportedType = func(t string) error { return fmt.Errorf("unsupported type: %s", t) }
	ErrorEmptyEnviron    = errors.New("environment can't be empty on strict mode")
)
