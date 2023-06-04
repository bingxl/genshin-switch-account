package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	infoBox *fyne.Container
	lib     *Lib
)

func setInfo(infos ...string) {
	for _, s := range infos {
		infoBox.Add(
			widget.NewLabel(s),
		)
	}
}

func selectGamePath(w fyne.Window) *fyne.Container {
	// gamePathLabel := widget.NewLabel(lib.config.GamePath)
	gamePathLabel := widget.NewLabel("点击选择游戏路径: ")

	btn := &widget.Button{
		Text:       lib.Config.Game,
		Importance: widget.HighImportance,
	}
	btn.OnTapped = func() {
		dialog.ShowFileOpen(func(url fyne.URIReadCloser, err error) {
			if err != nil {
				recover()
				setInfo(fmt.Sprintln(err))
				return
			}
			if url == nil {
				// 取消了选择或者没选到
				return
			}
			targetPath := url.URI().Path()

			btn.SetText(targetPath)
			lib.Config.GamePath = filepath.Dir(targetPath)
			lib.Config.Game = targetPath
			lib.SaveConfig()
		}, w)
	}

	return container.NewHBox(
		gamePathLabel,
		btn,
	)
}

func regList() *fyne.Container {

	lists := container.NewHBox()
	clickHandle := func(reg string) {
		lib.ChangeAccount(reg)
	}
	if len(lib.Regs) == 0 {
		//
		lists.Add(widget.NewLabel("请先导出注册表到" + filepath.Join(lib.CurrentPath, "reg")))
		return lists
	}

	for _, reg := range lib.Regs {
		btn := widget.NewButton(filepath.Base(reg), func() {
			clickHandle(reg)
		})
		btn.Importance = widget.HighImportance
		lists.Add(btn)
	}
	return lists
}

func main() {
	a := app.New()
	w := a.NewWindow("原神账号切换")
	w.Resize(fyne.NewSize(500, 500))

	infoBox = container.NewVBox()

	lib = &Lib{
		Log: setInfo,
	}

	lib.Init()

	w.SetContent(
		container.NewVBox(
			selectGamePath(w),
			widget.NewLabel("点击下方按钮切换到对应的账号"),
			regList(),
			// widget.NewButton("开始游戏", lib.StartGame),
			infoBox,
		))
	w.Show()
	a.Run()
}
