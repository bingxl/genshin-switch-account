package main

import (
	"context"
	"genshin/backend"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	lib *backend.Lib
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.init()
}

func (a *App) init() {
	a.lib = backend.NewLib()
	a.lib.Log = func(infos ...any) {
		a.emit("log", infos...)
	}
	a.lib.Init()
}

// 触发事件
func (a *App) emit(event string, params ...any) {
	runtime.EventsEmit(a.ctx, event, params...)
}

// 选择游戏可执行文件
func (a *App) SelectGameFile() bool {
	executeFilter := []runtime.FileFilter{
		{DisplayName: "execute Files *.exe", Pattern: "*.exe"},
	}

	gameFile, err := runtime.OpenFileDialog(a.ctx,
		runtime.OpenDialogOptions{
			Filters: executeFilter,
		})
	if err == nil {
		a.lib.SetGameFile(gameFile)
	}

	return err == nil
}

// 获取游戏可执行文件
func (a *App) GetGameFile() string {
	return a.lib.Config.Game
}

// 获取账号注册表文件
func (a *App) GetRegs() []string {
	a.lib.ReadRegs()
	return a.lib.Regs
}

// 将当前注册表内容导出到注册表文件
func (a *App) ExportReg(regName string) {
	if regName == "" {
		return

	}
	a.lib.Export(regName)

}

// 切换账号
func (a *App) ImportReg(regName string) {
	if regName == "" {
		return
	}
	a.lib.ChangeAccount(regName)
}

// 启动游戏
func (a *App) StartGame() {
	a.lib.StartGame()
}

// 启动 GIS
func (a *App) StartGis(gisPath string) {
	if gisPath == "" {
		return
	}
	a.lib.StartGis(gisPath)
}
