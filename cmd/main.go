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

	fmt.Println("索引\t", "账号")
	for i, v := range GetRegs() {
		fmt.Println(i, "\t", v)
	}
	fmt.Println("输入操作加索引， i开头导入账号， e开头导出账号， q退出")
	fmt.Println("eg: i0 导入第一个账号，e0 导出第一个账号")

	var input string

	for {
		fmt.Scanln(&input)
		if input == "q" || input == "exit" {
			break
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
