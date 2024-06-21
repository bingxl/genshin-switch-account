// 编译为命令行模式

package main

import (
	"fmt"
	"genshin/backend"
	"strconv"
)

var lib = backend.NewLib()

func main() {
	lib.Init()

	fmt.Println("索引  账号")
	for i, v := range GetRegs() {
		fmt.Println(i, "  ", v)
	}
	fmt.Println("输入操作后按回车\n i+索引号  导入账号 eg: i0 \n e+索引号  导出账号 eg: e1 \n q  退出 \n s  启动游戏")

	var input string

	for {
		fmt.Scanln(&input)
		if input == "q" || input == "exit" {
			lib.Close()
			break
		}
		if input == "s" {
			lib.StartGame()
			continue
		}
		operation := input[0]
		accountIndex, err := strconv.Atoi(input[1:])
		if err != nil || accountIndex >= len(GetRegs()) {
			fmt.Println("无效索引")
			continue
		}
		switch operation {
		case 'i':
			ImportReg(GetRegs()[accountIndex])
		case 'e':
			ExportReg(GetRegs()[accountIndex])
		default:
			fmt.Println("无效操作")
		}

	}
}

// 获取账号注册表文件
func GetRegs() []string {
	lib.ReadRegs()
	return lib.Regs
}

// 将当前注册表内容导出到注册表文件
func ExportReg(regName string) {
	if regName == "" {
		return
	}
	lib.Export(regName)
}

// 切换账号
func ImportReg(regName string) {
	if regName == "" {
		return
	}
	lib.ChangeAccount(regName)
}
