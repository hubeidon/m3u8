//go:build linux

package initial

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

var Aria2cPid *os.Process

func StartAria2c() {
	cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-all")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		panic("aria2c 启动失败")
	}
	Aria2cPid = cmd.Process
	log.Fatal(cmd.Wait())
}

func StopAria2c() {
	log.Error(syscall.Kill(-Aria2cPid.Pid, syscall.SIGKILL))
	log.Error(os.RemoveAll("data"))
}
