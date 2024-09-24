package main

import (
	"context"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	hook "github.com/robotn/gohook"
	"strconv"
	"time"
)

func main() {
	var mw *walk.MainWindow
	var te *walk.TextEdit
	var ti *walk.LineEdit

	// 创建主窗口
	err := MainWindow{
		AssignTo: &mw,
		Title:    "超级粘贴",
		Size:     Size{Width: 450, Height: 500},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "软件运行会自动监听剪切板\n也可以将内容复制到下面的文本框\n将光标放到需要复制的地方同时按下空格\n按下esc停止输入",
				Font: Font{
					Family:    "Microsoft YaHei",
					PointSize: 14,
					Bold:      true,
				},
			},
			Composite{
				Layout: HBox{
					Spacing:     0,
					Margins:     Margins{0, 0, 0, 0},
					MarginsZero: false,
					SpacingZero: true,
				},
				Children: []Widget{
					Label{
						Text: "每个字间隔：",
					},
					LineEdit{
						AssignTo: &ti,
						MinSize:  Size{50, 0},
						MaxSize:  Size{50, 0},
					},
					Label{
						Text: "毫秒（1s=1000ms）",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "粘贴",
						OnClicked: func() {
							text, err := walk.Clipboard().Text()
							if err != nil {
								te.SetText("")
							} else {
								te.SetText(text)
							}
						},
					},
					PushButton{
						Text: "清空",
						OnClicked: func() {
							te.SetText("")
						},
					},
				},
			},
			TextEdit{
				AssignTo: &te,
				VScroll:  true,
				ReadOnly: false,
				Font: Font{
					Family:    "Microsoft YaHei",
					PointSize: 8,
				},
			},
		},
	}.Create()

	if err != nil {
		panic(err)
	}
	s, _ := walk.Clipboard().Text()
	te.SetText(s)

	ti.SetText("2")
	var flags bool
	var ctx context.Context
	var cancel func()

	go func() {
		hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
			if cancel != nil {
				cancel()
			}
		})

		hook.Register(hook.KeyDown, []string{"space"}, func(e hook.Event) {
			flags = true

			str := te.Text()
			tim := ti.Text()
			timInt, _ := strconv.Atoi(tim)

			ctx, cancel = context.WithCancel(context.Background())

			go func() {
				for _, c := range str {
					select {
					case <-ctx.Done():
						return
					default:
						if !flags {
							break
						}
						robotgo.TypeStr(string(c))
						time.Sleep(time.Duration(timInt) * time.Millisecond)
					}
				}
				flags = false
			}()
		})

		ss := hook.Start()
		<-hook.Process(ss)
	}()

	// 运行窗口
	mw.Run()
}
