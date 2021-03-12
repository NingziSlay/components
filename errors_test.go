package pkg

import (
	"github.com/pkg/errors"
	"testing"
)

func TestConfigError_Error(t *testing.T) {
	var err = newE("some error")
	if err.Error() != "some error" {
		t.Fatalf("unexpected result: %s", err.Error())
	}

	err = wrapE("hello", errors.New("some error"))
	if err.Error() != errors.Wrap(err.err, err.msg).Error() {
		t.Fatalf("unexpected result: %s", err.Error())
	}
}
