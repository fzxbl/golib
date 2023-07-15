package iutil

import (
	"testing"
)

func Test_BlockIfExpired(t *testing.T) {
	BlockIfExpired(2023, 06, 26, 14)
}
