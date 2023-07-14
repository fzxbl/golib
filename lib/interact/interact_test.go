package interact

import (
	"fmt"
	"testing"
)

func Test_DirectorySelect(t *testing.T) {
	d, e := DirectorySelect("asdf")
	fmt.Println(d, e)
}

func Test_BlockUntilSignal(t *testing.T) {
	BlockOnSignal()
}

func Test_FileSelect(t *testing.T) {
	d, e := FileSelect("adsf", "")
	fmt.Println(d, e)
}
