package backend

import (
	"os"
	"path"
)

// 获取当前用户主目录
func getHomePath() (home string, err error) {
	home, err = os.UserHomeDir()
	if err != nil {
		home = "./"
	}
	return
}

// 获取项目目录, 主目录+项目名, 不存在则创建目录
func getProjectStorePath(projectName string) string {
	homeDir, _ := getHomePath()
	projDir := path.Join(homeDir, projectName)
	if !isExists(projDir) {
		os.Mkdir(projDir, 0750)
	}
	return projDir
}

func isExists(filenameOrDir string) bool {
	_, err := os.Stat(filenameOrDir)
	return os.IsExist(err)
}
