package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// 配置文件格式
type Config struct {
	LaunchPath string `json:"launchPath"`
	GamePath   string `json:"gamePath"`
}

var config, _ = readConfig()

var exePath, _ = os.Executable()
var currentPath = filepath.Dir(exePath)

// 读取配置文件
func readConfig() (*Config, error) {
	// io/ioutil 包已弃用，不能再用ioutil
	bytes, err := os.ReadFile(filepath.Join(currentPath, "./config.json"))
	if err != nil {
		fmt.Printf("读取配置文件出错 %+v", err)
		return nil, err
	}
	result := &Config{}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		fmt.Println("解码json数据时出错")
		return nil, err
	}
	return result, nil
}

// 获取注册表列表
func getRegs() ([]string, error) {
	matchs, err := filepath.Glob(filepath.Join(currentPath, "./reg/*.reg"))
	if err == nil {
		return matchs, nil
	}
	return nil, err
}

// 获取用户输入
func scanner(read *bufio.Reader, regs []string) {

	consoleStr := "\n ------请输入数字选择需要导入的注册表------"
	for index, reg := range regs {
		consoleStr += fmt.Sprintf("\n%d  %s", index, reg)
	}
	fmt.Println(consoleStr)

	var index int
	_, err := fmt.Scanln(&index)
	// fmt.Println("获得的数值：", index)
	// 输入数字不正确时让重输入
	if err != nil || index >= len(regs) || index < 0 {
		//
		fmt.Println("请输入正确的数字")
		scanner(read, regs)
		return
	}
	handle(regs[index])

}

// 更新游戏配置文件
func serverConfig(serverName string) {
	var modify [3]string
	switch serverName {
	case "bilibili":
		modify = [3]string{"14", "bilibili", "0"}
		// 将B服专用SDK复制到游戏目录下
		cpBiliBiliSDK()
	case "mihoyo":
		modify = [3]string{"1", "mihoyo", "1"}

	default:
		fmt.Println("暂不支持的server", serverName)
		return
	}

	changes := map[string]string{
		"channel":     modify[0],
		"cps":         modify[1],
		"sub_channel": modify[2],
	}

	gameConfigPath := filepath.Join(config.GamePath, "config.ini")
	// 读取 INI 文件
	cfg, err := ini.Load(gameConfigPath)
	if err != nil {
		fmt.Printf("无法读取 INI 文件：%v\n", err)
		return
	}

	// 获取或设置配置项的值
	section := cfg.Section("General")
	for k, v := range changes {
		section.Key(k).SetValue(v)
	}

	// 保存 INI 文件
	err = cfg.SaveTo(gameConfigPath)
	if err != nil {
		fmt.Printf("无法保存 INI 文件：%v\n", err)
		return
	}

	fmt.Println(gameConfigPath, "已更新并保存成功")

}

// 流程处理
func handle(reg string) {
	server := strings.Split(filepath.Base(reg), "-")[0]
	// TODO 依据 server 将 config 写入文件系统
	serverConfig(server)
	// 执行注册表导入
	importreg := exec.Command("reg", "import", reg)
	if err := importreg.Run(); err != nil {
		fmt.Println("导入注册表失败", err)
	} else {
		fmt.Println("导入注册表成功", reg)
	}

}

// 将 B 服SDK复制到指定位置
func cpBiliBiliSDK() {
	// 将sdk移动到对应位置， 此sdk为b服专有
	sourcePath := filepath.Join(currentPath, "./source/PCGameSDK.dll")
	targetPath := filepath.Join(config.GamePath, "YuanShen_Data", "Plugins", "PCGameSDK.dll")
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		// 读取错误处理
		fmt.Println("读取./source/PCGameSDK.dll SDK文件失败", err)
	}
	err = os.WriteFile(targetPath, content, 0755)
	if err != nil {
		// 写文件出错处理
		fmt.Println("SDK 写入游戏目录失败", err)
	}
	fmt.Println(targetPath)

}

func main() {

	// fmt.Printf("获取到的配置文件为： %+v", config)

	regs, _ := getRegs()

	read := bufio.NewReader(os.Stdin)

	for {
		// 循环读取用户输入
		scanner(read, regs)
	}
}
