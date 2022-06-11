package initial

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
)

var m3u8 *colly.Collector

func init() {
	c := colly.NewCollector()
	c.Async = true
	extensions.RandomUserAgent(c)

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Minute,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

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

	c.OnRequest(func(r *colly.Request) {
		log.Debug(r.URL.RequestURI())
	})

	// 解析m3u8文件
	c.OnResponse(func(r *colly.Response) {
		if getUrlFileFormat(r.Request.URL.String()) != "m3u8" {
			return
		}

		// 解析m3u8 文件
		m3u8Reader := bytes.NewReader(r.Body)
		buf := bufio.NewReader(m3u8Reader)
		for {
			readString, err := buf.ReadString('\n')
			if strings.HasPrefix(readString, "https://") {
				readString = strings.TrimSpace(readString)
				c.Visit(readString)
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

		Mkdir("./data/"+uri[39:55], os.ModePerm)

		i := strings.LastIndex(uri, "/")
		if err := r.Save(fmt.Sprintf("./data/%s/%s", uri[39:55], uri[i+1:])); err != nil {
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

//
func getUrlFileFormat(uri string) string {
	i := strings.Index(uri, "?")
	if i < 0 {
		i2 := strings.LastIndex(uri, ".")
		return uri[i2+1:]
	}
	i2 := strings.LastIndex(uri[:i], ".")
	if i < 0 {
		return ""
	}
	return uri[i2+1 : i]
}

var m = make(map[string]struct{})

func Mkdir(name string, perm os.FileMode) {
	if _, ok := m[name]; !ok {
		os.MkdirAll(name, perm)
		m[name] = struct{}{}
	}
}
