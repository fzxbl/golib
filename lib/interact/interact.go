package interact

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sqweek/dialog"
)

func ListenInput() (input string, err error) {
	inputReader := bufio.NewReader(os.Stdin)
	input, err = inputReader.ReadString('\n')
	if err != nil {
		return
	}
	input = strings.TrimSpace(input)
	return
}

// FileSelect 调用系统窗口进行单文件选择，main函数开始执行runtime.LockOSThread(),该函数调用完毕后调用runtime.UnlockOSThread()
func FileSelect(title, desc string, extensions ...string) (filePath string, err error) {
	filePath, err = dialog.File().Filter(desc, extensions...).Title(title).Load()
	return
}

func DirectorySelect(title string) (dir string, err error) {
	dir, err = dialog.Directory().Title(title).Browse()
	return
}

func BlockOnSignal() {
	log.Print("按 Ctrl-C 终止...")
	chSignal := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chSignal
	log.Print("程序1s后退出...")
	time.Sleep(time.Second)
	os.Exit(0)
}
