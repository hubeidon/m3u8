package initial

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
)

var (
	m3u8                    *colly.Collector
	tsFileNum               = make(map[string]int)
	notificationToStartWork = make(chan string, 3)
)

func init() {
	c := colly.NewCollector()
	c.Async = true

	extensions.RandomUserAgent(c)

	c.SetRequestTimeout(time.Minute)

	if err := c.Limit(&colly.LimitRule{
		DomainRegexp: `z.weilekangnet\.com`,
		RandomDelay:  500 * time.Millisecond,
		Parallelism:  5,
	}); err != nil {
		log.Fatal(err)
	}

	c.OnError(func(r *colly.Response, err error) {
		log.Errorf("url : %s , err : %v", r.Request.URL.String(), err)
	})

	// 解析m3u8文件
	c.OnResponse(func(r *colly.Response) {
		if getUrlFileFormat(r.Request.URL.String()) != "m3u8" {
			return
		}

		// 解析m3u8 文件
		m3u8Reader := bytes.NewReader(r.Body)
		buf := bufio.NewReader(m3u8Reader)
		var l bool
		for {
			readString, err := buf.ReadString('\n')
			if strings.HasPrefix(readString, "https://") {
				readString = strings.TrimSpace(readString)
				c.Visit(readString)
				tsFileNum[readString[56:72]]++

				if !l {
					notificationToStartWork <- readString[56:72]
					l = true
				}
			}

			if err != nil {
				return
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		uri := r.Request.URL.String()
		if getUrlFileFormat(uri) != "ts" {
			return
		}

		dirName := uri[56:72]

		Mkdir("./data/"+dirName, os.ModePerm)

		i := strings.LastIndex(uri, "/")
		if err := r.Save(fmt.Sprintf("./data/%s/%s", dirName, uri[i+1:])); err != nil {
			log.Error(err)
		}

	})

	m3u8 = c
}

func Down(uri string) {
	uri = strings.TrimSpace(uri)
	if err := m3u8.Visit(uri); err != nil {
		log.Errorf("%s,err: %v", uri, err)
	}
}

// getUrlFileFormat 返回url文件类型
// https://xxx.html -> html
// https://xxx -> ""
// https://xxx.ts -> ts
func getUrlFileFormat(uri string) string {
	i := strings.Index(uri, "?")
	if i < 0 {
		i2 := strings.LastIndex(uri, ".")
		return uri[i2+1:]
	}
	i2 := strings.LastIndex(uri[:i], ".")
	if i2 < 0 {
		return ""
	}
	return uri[i2+1 : i]
}

var m = make(map[string]struct{})

// Mkdir 创建文件
func Mkdir(name string, perm os.FileMode) {
	if _, ok := m[name]; !ok {
		os.MkdirAll(name, perm)
		m[name] = struct{}{}
	}
}

// CompositeVideo 等待下载器下载完毕，后合成视频
func CompositeVideo() {
	for dir := range notificationToStartWork {
		m3u8.Wait()

		workPath := fmt.Sprintf("./data/%s", dir)

		dirs, err := os.ReadDir(workPath)
		if err != nil {
			log.Error(err)
			continue
		}

		if len(dirs) != tsFileNum[dir] {
			log.Error("m3u8中的uri数量和实际下载完毕的不符")
			continue
		}

		var buf bytes.Buffer

		for i := 0; i < tsFileNum[dir]; i++ {
			name := fmt.Sprintf("out%d.ts", i)
			ts, _ := os.ReadFile(workPath + "/" + name)
			buf.Write(ts)
		}

		if f, err := os.OpenFile(dir+".mp4", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
			log.Error(err)
		} else {
			f.Write(buf.Bytes())
			f.Close()
			log.Infof("downloaded %s", dir)
		}
	}
}
