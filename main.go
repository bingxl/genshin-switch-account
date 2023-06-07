package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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

// 选择游戏路径
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
				setInfo(err.Error())
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

// 导入导出注册表
func regList() *fyne.Container {
	// 导出的注册表存放路径
	regPath := filepath.Join(lib.CurrentPath, "reg")

	// 当前选中的账号文件名
	var selectReg = ""
	// 导出当前注册表时所用的文件名
	var exportname = binding.NewString()

	radios := &widget.RadioGroup{}

	for _, reg := range lib.Regs {
		radios.Append(filepath.Base(reg))
	}

	importEntry := widget.NewEntryWithData(exportname)
	importEntry.PlaceHolder = "输入导出的文件名，以 b 或 g 开头"
	inputvalidator := func(s string) error {
		if s == "" {
			// setInfo("导出文件名不能为空")
			return fmt.Errorf("导出文件名不能为空")
		}
		if s[0] != 'b' && s[0] != 'g' {
			// setInfo("导出文件名必须以 b 或 g 开头，b为b服，g为官服")
			return fmt.Errorf("导出文件名必须以 b 或 g 开头，b为b服，g为官服")
		}

		return nil
	}
	importEntry.Validator = inputvalidator

	radios.OnChanged = func(reg string) {
		if len(reg) != 0 {
			selectReg = filepath.Join(regPath, reg)
			exportname.Set(reg)

		}
	}

	bts := container.New(
		layout.NewAdaptiveGridLayout(2),
		&widget.Button{
			Text: "切换账号",
			OnTapped: func() {
				if len(selectReg) != 0 {
					lib.ChangeAccount(selectReg)
				} else {
					setInfo("请先选中一个账号")
				}
			},
		},
		&widget.Button{
			Text: "导出当前注册表",
			OnTapped: func() {
				reg, err := exportname.Get()
				info := ""
				if err != nil {
					info = err.Error()
				} else if err = inputvalidator(reg); err != nil {
					info = err.Error()
				} else {
					lib.Export(filepath.Join(regPath, reg))
					return
				}
				setInfo(info)

			},
		},
	)

	return container.NewVBox(radios, importEntry, bts)
}

func main() {
	a := app.New()
	w := a.NewWindow("原神账号切换")
	w.Resize(fyne.NewSize(500, 500))

	infoBox = container.New(layout.NewVBoxLayout())
	infoScroll := container.NewVScroll(infoBox)
	infoScroll.SetMinSize(fyne.NewSize(0, 300))

	logArea := container.New(
		layout.NewAdaptiveGridLayout(2),
		widget.NewLabel("日志区域:"),
		widget.NewButton("清除日志", func() {
			infoBox.RemoveAll()
		}),
	)

	lib = &Lib{
		Log: setInfo,
	}

	lib.Init()

	w.SetContent(
		container.NewVBox(
			selectGamePath(w),
			widget.NewLabel("点击下方按钮切换到对应的账号"),
			regList(),
			widget.NewButton("开始游戏", lib.StartGame),
			logArea,
			infoScroll,
		))
	w.Show()
	a.Run()
}
