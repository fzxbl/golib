package ilog

import (
	"testing"
)

func Test_InitLogger(t *testing.T) {
	logger := MustInitFromFile("testdata/logger.toml")
	logger.Info("asdf a")
}
