package main

import (
	"os"

	"github.com/flopp/go-findfont"
)

func init() {
	// 设置字体
	// 字体名称在 C:\\windows\fonts\ 下查找，找文件名
	fontfilename := "STFANGSO.TTF"
	fontPath, err := findfont.Find(fontfilename)

	if err != nil {
		panic(err)
	}
	os.Setenv("FYNE_FONT", fontPath)
}
