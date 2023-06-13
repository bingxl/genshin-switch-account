package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"gopkg.in/ini.v1"
)

// 配置文件格式
type ConfigT struct {
	LaunchPath string `json:"launch_path"`
	GamePath   string `json:"game_path"`
	Game       string `json:"game"`
}

type Lib struct {
	Config      *ConfigT
	Regs        []string
	CurrentPath string
	Log         func(...string)
	regKey      string
}

func (lib *Lib) Init() {
	var exePath, _ = os.Executable()
	lib.CurrentPath = filepath.Dir(exePath)
	lib.regKey = "HKCU\\Software\\miHoYo\\原神"

	lib.readConfig()
	lib.readRegs()

}

// 读取配置文件
func (lib *Lib) readConfig() {
	// io/ioutil 包已弃用，不能再用ioutil
	bytes, err := os.ReadFile(filepath.Join(lib.CurrentPath, "./config.json"))
	if err != nil {
		lib.logInfo("读取配置文件出错 %+v", err)
		return
	}
	result := &ConfigT{}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		lib.logInfo("解码json数据时出错")
		return
	}
	lib.Config = result
}

// 写入配置文件
func (lib *Lib) SaveConfig() {
	lib.logInfo(lib.Config)

	data, _ := json.Marshal(lib.Config)
	err := os.WriteFile(
		filepath.Join(lib.CurrentPath, "./config.json"),
		data,
		0755,
	)
	if err != nil {
		lib.logInfo("写入配置文件出错", err)
		return
	}
	lib.logInfo("成功写入配置文件")
}

// 获取注册表列表
func (lib *Lib) readRegs() {
	matchs, err := filepath.Glob(filepath.Join(lib.CurrentPath, "./reg/*.reg"))
	if err != nil {
		lib.logInfo("读取注册表文件出错", err)
	}
	lib.Regs = matchs
}

// 导出注册表到
func (lib *Lib) Export(file string) {
	err := lib.RunCommand("reg", "export", lib.regKey, file, "/y")
	if err != nil {
		lib.logInfo("导出注册表失败", err)
	} else {
		lib.logInfo("导出注册表成功")
	}
}

// 更新游戏配置文件
func (lib *Lib) serverConfig(serverName byte) {
	var modify [3]string
	switch serverName {
	case 'b':
		modify = [3]string{"14", "bilibili", "0"}
	case 'g':
		modify = [3]string{"1", "mihoyo", "1"}

	default:
		lib.logInfo("暂不支持的server", serverName)
		return
	}
	// 将B服专用SDK复制到游戏目录下
	lib.cpBiliBiliSDK(serverName != 'b')

	changes := map[string]string{
		"channel":     modify[0],
		"cps":         modify[1],
		"sub_channel": modify[2],
	}

	gameConfigPath := filepath.Join(lib.Config.GamePath, "config.ini")
	// 读取 INI 文件
	cfg, err := ini.Load(gameConfigPath)
	if err != nil {
		lib.logInfo("无法读取 INI 文件：%v\n", err)
		return
	}

	// 获取或设置配置项的值
	section := cfg.Section("General")

	// 读取的配置和将要写入的配置是否相同
	hasDiff := false
	for k, v := range changes {
		if section.Key(k).String() != v {
			hasDiff = true
		}
		section.Key(k).SetValue(v)
	}

	// 配置文件与将要写入的内容没有不同，则直接返回
	if !hasDiff {
		return
	}

	// 保存 INI 文件
	err = cfg.SaveTo(gameConfigPath)
	if err != nil {
		lib.logInfo("无法保存 INI 文件：%v\n", err)
		return
	}

	lib.logInfo(gameConfigPath, "已更新并保存成功")
}

// 流程处理
func (lib *Lib) ChangeAccount(reg string) {
	lib.logInfo("changeAccount 接收到：", reg)
	server := filepath.Base(reg)[0]
	lib.serverConfig(server)

	// 执行注册表导入
	err := lib.RunCommand("reg", "import", reg)
	if err != nil {
		lib.logInfo("导入注册表失败", err)
	} else {
		lib.logInfo("导入注册表成功", reg)
	}

}

// 目标为 B 服时如果游戏目录下没有SDK则拷贝SDK，为官服时若游戏目录下有 SDK 则删除
func (lib *Lib) cpBiliBiliSDK(remove bool) {
	// 将sdk移动到对应位置， 此sdk为b服专有
	sourcePath := filepath.Join(lib.CurrentPath, "./source/PCGameSDK.dll")
	targetPath := filepath.Join(lib.Config.GamePath, "YuanShen_Data", "Plugins", "PCGameSDK.dll")

	// 判断文件是否存在
	_, err := os.Stat(targetPath)
	sdkHasExist := err == nil

	if remove && sdkHasExist {
		os.Remove(targetPath)
		return
	}

	// 需要copy SDK 且 SDK 不存在
	if !remove && !sdkHasExist {
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			// 读取错误处理
			lib.logInfo("读取./source/PCGameSDK.dll SDK文件失败", err)
		}
		err = os.WriteFile(targetPath, content, 0755)
		if err != nil {
			// 写文件出错处理
			lib.logInfo("SDK 写入游戏目录失败", err)
		}
	}

	lib.logInfo(targetPath)

}

// 将数据显示到界面中
func (lib *Lib) logInfo(args ...interface{}) {
	s := ""
	for _, arg := range args {
		s += fmt.Sprintf("%+v", arg)
	}
	if lib.Log != nil {
		lib.Log(s)
	}
	fmt.Println(s)

}

// 开始游戏
func (lib *Lib) StartGame() {

	cmd := exec.Command(lib.Config.Game)
	if err := cmd.Start(); err != nil {
		lib.logInfo("游戏启动失败", err)
		return
	}
	// cmd.Wait()
	lib.logInfo("游戏启动中")

}

// 运行命令
func (lib *Lib) RunCommand(args ...string) error {

	cmd := exec.Command(args[0], args[1:]...)

	// exe执行时会启动一个终端，不隐藏 Window 时会有终端闪现
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()

}
