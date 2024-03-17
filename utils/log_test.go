package utils

import (
	"errors"
	"testing"
)

func TestLogError(t *testing.T) {
	LogError("gzf 666", errors.New("你错了"))
}
