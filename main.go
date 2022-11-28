package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gitee.com/don178/m3u8/global"
	"gitee.com/don178/m3u8/initial"
	"github.com/gocolly/colly/v2"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// isExistsHost 判断ts路径是否需要添加前缀
func isExistsHost(in []byte) bool {
	var buf = bytes.NewBuffer(in)
	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		if bytes.HasPrefix(line, []byte("http")) {
			return true
		}
	}
	return false
}

// parseCompleteM3u8
func parseCompleteM3u8(in []byte) []string {

	var res = make([]string, 0, 50)
	var buf = bytes.NewBuffer(in)

	for {
		line, err := buf.ReadBytes('\n')

		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		if err == io.EOF {
			break
		}

		if bytes.HasPrefix(line, []byte("http")) {
			res = append(res, string(bytes.TrimSpace(line)))
		}
	}

	return res
}

// parseNotCompleteM3u8
func parseNotCompleteM3u8(in []byte, host string) []string {
	var res = make([]string, 0, 50)
	var buf = bytes.NewBuffer(in)

	for {
		line, err := buf.ReadBytes('\n')
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		if err == io.EOF {
			break
		}
		res = append(res, fmt.Sprintf("%s/%s", host, bytes.TrimSpace(line)))
	}

	return res
}

const (
	HTTP  = "HTTP"
	LOCAL = "LOCAL"
)

type m3u8 struct {
	// 文件内容
	body bytes.Buffer
	// tsBody 当需要视频需要解密时,存放未解密前的数据
	tsBody bytes.Buffer
	// 解析出来的url地址
	urls []string
	// 来源类型
	sType string
	// 来源
	source string
	// 文件名
	name string
	// 下载器
	*colly.Collector
	// 前缀
	prefix string
	// 解密算法
	initial.Decryption
}

func NewM3u8ByAddress(address global.Address) *m3u8 {
	return NewM3u8(address.Path, address.Prefix, address.Fname)
}

func NewM3u8(path, prefix, fname string) *m3u8 {
	m := new(m3u8)

	m.sType = isHttpOrLocal(path)
	m.source = path
	m.prefix = prefix
	m.name = fname

	m.Collector = func() *colly.Collector {
		c := colly.NewCollector()
		c.MaxBodySize = -1
		c.UserAgent = global.Cfg.UserAgent
		c.SetRequestTimeout(global.Cfg.RequestTimeout * time.Second)
		return c
	}()

	return m
}

// isHttpOrLocal 判断是网络地址还是本地地址
func isHttpOrLocal(path string) string {
	if strings.HasPrefix(path, "http") {
		return HTTP
	}

	if f, err := os.Open(path); err != nil {
		global.Log.Fatal(path, zap.Error(err))
	} else {
		f.Close()
	}
	return LOCAL
}

func (m *m3u8) getTsOnHttp() error {
	res, err := http.Get(m.source)
	if err != nil {
		return err
	}

	n, err := m.body.ReadFrom(res.Body)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("body为空 err:%w", err)
	}
	return nil
}
func (m *m3u8) getTsOnLocal() error {
	f, err := os.Open(m.source)
	if err != nil {
		return err
	}
	n, err := m.body.ReadFrom(f)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("body为空 err:%w", err)
	}
	return nil
}
func (m *m3u8) getM3u8() {
	switch m.sType {
	case HTTP:
		m.getTsOnHttp()
	case LOCAL:
		m.getTsOnLocal()
	default:
		global.Log.Fatal("stype error", zap.String("sType", m.sType))
	}
}

func (m *m3u8) parseUrls() {
	if isExistsHost(m.body.Bytes()) {
		m.urls = parseCompleteM3u8(m.body.Bytes())
	} else {
		if m.prefix == "" {
			m.prefix = m.source[:strings.LastIndex(m.source, "/")]
		}
		m.urls = parseNotCompleteM3u8(m.body.Bytes(), m.prefix)
	}
	// 特殊情况
	// 1. m3u8内所以ts路径都是一样的, 前端通过Range获取不同ts碎片
	if m.urls[0] == m.urls[len(m.urls)-1] {
		m.urls = m.urls[:1]
	}
}

func (m *m3u8) fileName() {
	var n = -1
	for i, str := range m.source {
		if str == '/' {
			n = i
		} else if str == '?' {
			m.name = m.source[n+1 : i]
			return
		}
	}
	m.name = m.source[n+1:]
}

// completePath 返回文件完整地址
func (m *m3u8) completePath() string {
	return filepath.Join(global.Cfg.Dir, fmt.Sprint(m.getName(), global.Cfg.Ext))
}

// getName 获取文件名 (不包含配置文件ext扩展名)
func (m *m3u8) getName() string {
	if m.name == "" {
		m.fileName()
	}
	return m.name
}

// onColly 注册collyResponse
func (m *m3u8) onColly() {
	if m.isEncryption() {
		m.OnResponse(func(r *colly.Response) {
			global.Log.Debug("写入文件:", zap.String("m3u8", m.getName()), zap.String("ts", r.FileName()))
			n, err := m.tsBody.Write(r.Body)
			if fmt.Sprint(n) != r.Headers.Get("content-length") {
				global.Log.Error("write m.tsBody error", zap.Error(err))
			}
		})
	} else {
		m.OnResponse(func(r *colly.Response) {
			global.Log.Debug("写入文件:", zap.String("m3u8", m.getName()), zap.String("ts", r.FileName()))
			f, err := os.OpenFile(m.completePath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
			if err != nil {
				global.Log.Error("save file ", zap.Error(err))
			}
			f.Write(r.Body)
			f.Close()
		})
	}
}

// isEncryption 判断是否有加密, 如果有会初始化 m3u8.Decryption(解密处理程序)
func (m *m3u8) isEncryption() bool {
	for {
		line, err := m.tsBody.ReadBytes('\n')
		if err != nil {
			break
		}
		if bytes.HasPrefix(line, []byte("#EXT-X-KEY:METHOD=")) {
			break
		}
		if line[0] != '#' {
			return false
		}
	}

	for _, decryption := range initial.DecryptionList {
		if decryption.IsNeed(m.tsBody.Bytes()) {
			m.Decryption = decryption
			return true
		}
	}
	return false
}

func (m *m3u8) download() error {
	for _, url := range m.urls {
		if err := m.Visit(url); err != nil {
			return err
		}
	}
	return nil
}

// decrypt 进行视频解密后,保存到文件.
func (m *m3u8) decrypt() error {
	if m.Decryption == nil {
		return errors.New("解密处理程序为空")
	}

	dst, err := m.Decryption.Decrypt(m.tsBody.Bytes())
	if err != nil {
		return err
	}

	f, err := os.Create(m.completePath())
	if err != nil {
		return err
	}

	if n, err := f.Write(dst); n != len(dst) {
		global.Log.Error("加密后吸入文件错误", zap.Error(err))
	}
	return f.Close()
}

func (m *m3u8) Run() error {
	global.Log.Debug("开始下载:", zap.String("m3u8", m.getName()))
	m.getM3u8()
	m.parseUrls()
	m.onColly()
	err := m.download()
	if err != nil {
		return err
	}
	global.Log.Debug("下载完毕:", zap.String("m3u8", m.getName()))
	if m.Decryption != nil {
		global.Log.Debug("解密中:", zap.String("m3u8", m.getName()))
		return m.decrypt()
	}
	return nil
}

func main() {
	pool, err := ants.NewPool(global.Cfg.GoNum)
	if err != nil {
		global.Slog.Fatalln(err)
	}
	var wg sync.WaitGroup
	for _, address := range global.Cfg.Address {
		m := NewM3u8ByAddress(address)
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			if err := m.Run(); err != nil {
				global.Slog.Errorln(err)
			}
		})
	}
	wg.Wait()
}
