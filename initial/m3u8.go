package initial

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"m3u8/aria2"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	aria    *aria2.AriaEngine
	Mp4Chan = make(chan Mp4, 10)
)

type Mp4 struct {
	Dir string
	Out string
	Num uint
}

func Mkdir(path string) string {
	for {
		//判断文件是否存在
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			_ = os.MkdirAll(path, os.ModePerm)
			return path
		} else {
			// 是否是文件夹 , 存在添加后缀 -cp
			path += "-cp"
		}
	}
}

func JoinMp4(mp4 <-chan Mp4) {
	log.Info("启动 视频合成 goroutine")
	for v := range mp4 {
		// 等待任务下载完成
		log.Info("%#v 等待合成", mp4)
		for {
			if _, err := os.Stat(fmt.Sprintf("%s/out%d.ts", v.Dir, v.Num)); errors.Is(err, fs.ErrNotExist) {
				time.Sleep(time.Second)
				continue
			}
			break
		}
		log.Info("%#v 开始合成", mp4)

		files, err := os.ReadDir(v.Dir)
		if err != nil {
			continue
		}
		length := len(files)
		var buf bytes.Buffer

		for i := 0; i < length; i++ {
			name := fmt.Sprintf("out%d.ts", i)
			ts, _ := os.ReadFile(v.Dir + "/" + name)
			buf.Write(ts)
		}

		for {
			if _, err := os.Stat(v.Out); errors.Is(err, fs.ErrNotExist) {
				f, _ := os.Create(v.Out)
				log.Info(f.Write(buf.Bytes()))
				log.Info("写入成功,关闭文件:", f.Close())
				break
			} else {
				res := strings.Split(v.Out, ".")
				v.Out = res[0] + "-cp.mp4"
			}
		}
	}
}

type OneDown struct {
	// m3u8Path 文件路径 可能是本地文件或者http url
	m3u8Path string
	// oneType 本地文件:local, http url: http
	oneType string
	// m3u8FileName 文件名称
	m3u8FileName string
	// dir 文件存储路径（文件夹）
	dir string
	// m3u8Byte 文件内容
	m3u8Byte []byte
	// 视频切边数量
	num uint
}

func (one *OneDown) parseM3u8() {
	m3u8Reader := bytes.NewReader(one.m3u8Byte)
	buf := bufio.NewReader(m3u8Reader)
	n := 0
	for {
		readString, err := buf.ReadString('\n')
		if strings.HasPrefix(readString, "https://") {
			readString = strings.TrimSpace(readString)
			n++
			aria.Download(readString, one.dir, fmt.Sprintf("out%d.ts", n), "")
		}
		if err != nil {
			one.num = uint(n)
			return
		}
	}
}

func NewOneDown(m3u8Path string) *OneDown {
	_, fileName := path.Split(m3u8Path)
	dir, _ := filepath.Abs(fmt.Sprintf("./data/%s", fileName))
	dir = strings.ReplaceAll(dir, `\`, `/`)
	dir = Mkdir(dir)

	m3u8Byte, err := os.ReadFile(m3u8Path)
	if err != nil {
		return nil
	}
	return &OneDown{
		m3u8Path:     m3u8Path,
		m3u8FileName: fileName,
		dir:          dir,
		m3u8Byte:     m3u8Byte,
		oneType:      "local",
	}
}

func NewOneHttp(httpUrl string) *OneDown {
	parse, err := url.Parse(httpUrl)
	if err != nil {
		log.Error(err)
		runtime.Goexit()
	}
	_, m3u8FileName := path.Split(parse.Path)

	dir, _ := filepath.Abs(fmt.Sprintf("./data/%s", m3u8FileName))
	dir = strings.ReplaceAll(dir, `\`, `/`)
	dir = Mkdir(dir)

	code, body, _ := fasthttp.GetTimeout([]byte{}, httpUrl, time.Second*5)
	if code != 200 {
		log.Error(err)
		runtime.Goexit()
	}
	return &OneDown{
		m3u8Path:     httpUrl,
		m3u8FileName: m3u8FileName,
		dir:          dir,
		m3u8Byte:     body,
		oneType:      "http",
	}
}

func HttpOrLocal(text string) {
	if strings.HasPrefix(text, "http") {
		handler(NewOneHttp(text))
	} else if filepath.IsAbs(text) {
		handler(NewOneDown(text))
	}
}

func handler(one *OneDown) {
	one.parseM3u8()

	abs, err := filepath.Abs(one.m3u8FileName + ".mp4")
	abs = strings.ReplaceAll(abs, `\`, `/`)
	if err != nil {
		return
	}
	Mp4Chan <- Mp4{
		Dir: one.dir,
		Out: abs,
		Num: one.num,
	}
}

func Run() {
	// 启动日志
	InitLogger()

	// 启动aria2c rpc服务
	go StartAria2c()

	// 合成碎片视频
	go JoinMp4(Mp4Chan)

	// 循环链接aria2c rpc
	for i := 0; i < 10; i++ {
		aria = aria2.NewAriaEngine()

		if aria != nil {
			break
		}
		time.Sleep(time.Second / 2)
	}
}

// CloseWindow 关闭窗口后执行
// 删除data目录
func CloseWindow() {
	cmd := exec.Command("rm", "-rf", "data")
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	log.Fatalf("aria2c err : %v", cmd.Run().Error())
}
