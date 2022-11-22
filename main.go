package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gitee.com/don178/m3u8/global"
	_ "gitee.com/don178/m3u8/initial"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

var (
	Meta   sync.Map
	M3u8er = InitM3u8()
	Downer = InitDowner()
)

type MetaVlaue struct {
	// fileName文件名称
	fileName string
	// dir 文件保存路径
	dir string
}

func (m MetaVlaue) completePath() string {
	return fmt.Sprintf("%s/%s", m.dir, m.fileName)
}

func InitM3u8() *colly.Collector {
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {

		// 两种情况
		// 1. m3u8内所有文件地址都是一样的(没有区分地址)(虎牙)
		// 2. m3u8内对地址进行了区分
		if err := netAddressHandler(r.Request.URL.String(), r.Body); err != nil {
			global.Slog.Error(err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		global.Slog.Error(err)
	})

	return c
}
func InitDowner() *colly.Collector {
	c := colly.NewCollector()
	c.MaxBodySize = -1
	c.UserAgent = global.Viper.GetString("user-agent")
	c.SetRequestTimeout(time.Second * time.Duration(global.Viper.GetInt("request-timeout")))

	c.OnResponse(func(r *colly.Response) {
		key := fmt.Sprintf("%x", md5.Sum([]byte(r.Request.URL.String())))
		if val, ok := Meta.Load(key); ok {
			m := val.(MetaVlaue)
			r.Save(m.completePath())
		} else {
			global.Log.Error("Meta miss", zap.String("url", r.Request.URL.String()))
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		global.Slog.Error(err, r.Request.URL.String())
	})
	return c
}

// netAddressHandler 处理网络m3u8
func netAddressHandler(uri string, body []byte) error {
	return m3u8Handler(uri[:lastLeftSlash(uri)], urlFileName(uri), body)
}

// m3u8Handler m3u8 文件处理
// prefix 当文件内没有host时会自动在路径前添加prefix
// dir 下载视频存储目录(碎片ts) 当文件内的ts路径是同一个时，使用配置文件内的dir，该dir无效
func m3u8Handler(prefix, dir string, body []byte) error {
	// urls 存储完整的ts路径
	var urls []string
	if isExistsHost(body) {
		urls = parseCompleteM3u8(body)
	} else {
		urls = parseNotCompleteM3u8(body, prefix)
	}

	if len(urls) < 1 {
		return fmt.Errorf("没有解析出ts路径")
	}

	// 当文件内的ts路径是同一个时
	if urls[0] == urls[len(urls)-1] {
		key := fmt.Sprintf("%x", md5.Sum([]byte(urls[0])))
		Meta.Store(key, MetaVlaue{
			fileName: urlFileName(urls[0]),
			dir:      global.Viper.GetString("dir"),
		})
		if err := Downer.Visit(urls[0]); err != nil {
			return err
		}
	} else {
		// 当文件内的ts路径不同
		for _, url := range urls {
			key := fmt.Sprintf("%x", md5.Sum([]byte(url)))
			dir := filepath.Join(global.Viper.GetString("dir"), dir)
			os.Mkdir(dir, os.ModePerm)
			Meta.Store(key, MetaVlaue{
				fileName: urlFileName(url),
				dir:      dir,
			})
			if err := Downer.Visit(url); err != nil {
				return err
			}
		}
	}
	return nil
}

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
			res = append(res, fmt.Sprint(bytes.TrimSpace(line)))
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

func lastLeftSlash(url string) int {
	var n int = -1
	for i, str := range url {
		if str == '/' {
			n = i
		} else if str == '?' {
			break
		}
	}
	return n
}

func urlFileName(url string) string {
	var n = -1

	for i, str := range url {
		if str == '/' {
			n = i
		} else if str == '?' {
			return url[n+1 : i]
		}
	}
	return url[n+1:]
}

func localFileHandler(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	prefix := global.Viper.GetString("prefix." + urlFileName(path))
	if err = m3u8Handler(prefix, urlFileName(path), body); err != nil {
		return err
	}

	return nil
}

func main() {
	websites := global.Viper.GetStringSlice("website")
	for _, website := range websites {
		M3u8er.Visit(website)
	}

	local := global.Viper.GetStringSlice("local")
	for _, l := range local {
		if err := localFileHandler(l); err != nil {
			global.Slog.Error(err)
		}
	}
}
