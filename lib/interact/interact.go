package interact

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/sqweek/dialog"
)

// ListenInput 等待用户输入，并打印用户输入到标准输出
func ListenInput() (input string, err error) {
	inputReader := bufio.NewReader(os.Stdin)
	input, err = inputReader.ReadString('\n')
	if err != nil {
		return
	}
	input = strings.TrimSpace(input)
	fmt.Println("你的输入为:", input)
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

// NeedConfirm 打印引导文本，使程序停下来等待用户确认
func NeedConfirm() bool {
	if err := keyboard.Open(); err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("按Enter继续，ESC退出...")
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)
			return false
		}

		if key == keyboard.KeyEnter {
			// 程序开始
			keyboard.Close()
			break
		} else if key == keyboard.KeyEsc {
			keyboard.Close()
			os.Exit(0)
			return false
		}
	}

	return true
}
