package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/fzxbl/golib/lib/interact"
)

func main() {
	interact.BlockOnSignal()
	a := app.New()
	w := a.NewWindow("hello")

	//第一步是将要更新的小部件分配给变量。
	clock := widget.NewLabel("time")
	c := container.New(layout.NewGridLayout(5), widget.NewLabel("Hello World!"), clock)

	go func() {
		for range time.Tick(time.Second) {
			formatted := time.Now().Format("03:04:05")
			clock.SetText(formatted)
		}
	}()
	w.Resize(fyne.NewSize(1000, 700))
	w.SetContent(c)
	w.ShowAndRun()
	tidyUp()
}
func tidyUp() {
	fmt.Println("Exited")
}
