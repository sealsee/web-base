package fileUtils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// CreateMutiDir 调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func HomePath() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if runtime.GOOS == "windows" {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// GetCurrentPath 获取当前路径，比如：E:/abc/data/test
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}
