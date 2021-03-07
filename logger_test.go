package components

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestInitLogger(t *testing.T) {
	InitLogger(true)

	if &log == nil {
		t.Fatalf("log should have been initialized")
	}

	if log.GetLevel() != zerolog.DebugLevel {
		t.Fatalf("log's level should be Debug")
	}

	InitLogger(false)

	if &log == nil {
		t.Fatalf("log should have been initialized")
	}

	if log.GetLevel() != zerolog.InfoLevel {
		t.Fatalf("log's level should be Info")
	}
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()

	if &logger == nil {
		t.Fatalf("logger should not be nil")
	}
}

func TestWithLoggerLevel(t *testing.T) {
	log := WithLoggerLevel(zerolog.TraceLevel)
	if log.GetLevel() != zerolog.TraceLevel {
		t.Fatalf("log level should be Trace")
	}

	log = WithLoggerLevel(zerolog.FatalLevel)
	if log.GetLevel() != zerolog.FatalLevel{
		t.Fatalf("log level should be Fatal")
	}
}
