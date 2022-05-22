//go:build windows

package initial

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

func StartAria2c() {
	path := os.Getenv("path")

	abs, err := filepath.Abs("")
	if err != nil {
		log.Fatalf("获取绝对路径失败！err : %v", err)
	}

	err = os.Setenv("path", path+";"+abs)
	if err != nil {
		log.Fatalf("设置环境变量失败！err : %v", err)
	}

	cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-all")
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	log.Fatalf("aria2c err : %v", cmd.Run().Error())
}

func StopAria2c() {
	c := exec.Command("taskkill.exe", "/f", "/im", "aria2c.exe")
	log.Error("关闭下载器 : ", c.Start())
	log.Error("删除视频切片数据 : ", os.RemoveAll("data"))
}

// CloseWindow 关闭窗口后执行
// 删除data目录
// func CloseWindow() {
// 	cmd := exec.Command("rm", "-rf", "data")
// 	if runtime.GOOS == "windows" {
// 		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
// 	}
// 	log.Fatalf("aria2c err : %v", cmd.Run().Error())
// }