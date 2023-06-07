# 原神多账号切换
原神多账号切换

## 实现的功能
实现官服多个账号切换；B 服只能切换到最后一次登录的账号，不能做到B服多个账号之前的切换

## 使用方法

1. 到[release](https://github.com/bingxl/genshin-switch-account/releases) 页面下载 对应版本的`genshin-switch-account-Vxxx.zip`文件.     
解压到一个空目录中；
2. 运行游戏登录账号，然后将注册表导出到本项目的 reg 文件夹下，命名格式为 `m-xxxx.reg`或者 `b-xxxx.reg`，其中xxxx自由输入，建议和和自己的账号关联 `m`开头 表示官服账号， `b`开头为 B 服账号；
3. 右键点击 GenshinSwitch.exe 以管理员身份运行；（右键->属性->兼容性， 勾选以管理员身份运行此程序 后双击就可以以管理员身份运行）
4. 在界面里点击选择游戏路径，选择具体的游戏文件名；
5. 点击需要切换的账号；
6. 点击运行游戏。   


## 导出注册表方法
有多种方法导出注册表：
1. 在程序界面中 输入导出的文件名后点击 `导出当前注册表` 按钮
2. 使用 windows 提供的 register 工具。
    按下 Win 键，然后在输入框中输入 reg 回车，点击出现的注册表编辑器，然后找到 `HKEY_CURRENT_USER\Software\miHoYo\原神`， 右键点击原神，选择导出

3. 命令行导出。
    打开终端 执行 `reg export HKCU\\Software\\miHoYo\\原神  导出后的文件名 /y`

## build
安装 fyne 命令行工具 `go install fyne.io/fyne/v2/cmd/fyne@latest `

打包`fyne package -os windows`
