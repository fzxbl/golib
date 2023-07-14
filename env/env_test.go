package ienv

import (
	"fmt"
	"testing"

	"github.com/sqweek/dialog"
)

func TestEnv(t *testing.T) {
	fmt.Println(RootDir())
	fmt.Println(ConfDir())
}

func extractScientificNotation() {
	filePath, err := dialog.File().Filter("文本文件", "bdf").Load()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Selected file:", filePath)
}

func Test_extractScientificNotation(t *testing.T) {
	extractScientificNotation()

}
