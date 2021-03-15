package pkg

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestInitLogger(t *testing.T) {
	InitLogger(true)

	if &l == nil {
		t.Fatalf("l should have been initialized")
	}

	if l.GetLevel() != zerolog.DebugLevel {
		t.Fatalf("l's level should be Debug")
	}

	InitLogger(false)

	if &l == nil {
		t.Fatalf("l should have been initialized")
	}

	if l.GetLevel() != zerolog.InfoLevel {
		t.Fatalf("l's level should be Info")
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
		t.Fatalf("l level should be Trace")
	}

	log = WithLoggerLevel(zerolog.FatalLevel)
	if log.GetLevel() != zerolog.FatalLevel {
		t.Fatalf("l level should be Fatal")
	}
}
