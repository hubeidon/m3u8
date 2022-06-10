package initial

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

var m3u8 *colly.Collector

func init() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	// 解析m3u8文件
	c.OnResponse(func(r *colly.Response) {
		if getUrlFileFormat(r.Request.URL.String()) != "m3u8" {
			return
		}
		
	})

	c.OnResponse(func(r *colly.Response) {

	})

	m3u8 = c
}

func Down(url string) {
	m3u8.Visit(url)
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
