package backend

// 导入导出注册表的命令， 根据需要调整
var regCmd = []string{"flatpak", "run", "--env=WINEPREFIX=/home/bingxl/.wine", "org.winehq.Wine", "reg"}
