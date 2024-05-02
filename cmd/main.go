package main

import (
	"fmt"
	"genshin/backend"
	"strings"
)

var lib = backend.NewLib()

func main() {
	lib.Init()

	fmt.Println(strings.Join(GetRegs(), "\n"))
	fmt.Println("输入操作， i开头导入账号， e开头导出账号， q退出")

	var input string

	for {
		fmt.Scanln(&input)
		if input == "q" {
			break
		}
		if strings.HasPrefix(input, "i") {
			ImportReg(input[1:] + ".reg")
		} else if strings.HasPrefix(input, "e") {
			ExportReg(input[1:] + ".reg")
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
