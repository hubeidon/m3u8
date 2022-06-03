package aria2

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

//add.json.给出下载地址后发送到aria2
//2.设置请求参数，headers，referer等
//3.判断某个请求是否完成
//4.判断某些请求是否完成

// AriaEngine 调用aria2下载器
type AriaEngine struct {
	// http Request
	req *fasthttp.Request
	// http Response
	resp *fasthttp.Response
}

func NewAriaEngine() *AriaEngine {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod("POST")
	req.Header.SetRequestURI("http://localhost:6800/jsonrpc")
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Error(err)
		return nil
	}
	resp.Reset()
	return &AriaEngine{
		req:  req,
		resp: resp,
	}
}

func (aria *AriaEngine) CloseEngine() {
	fasthttp.ReleaseRequest(aria.req)
	fasthttp.ReleaseResponse(aria.resp)
}

func (aria *AriaEngine) Download(url, dir, out, referer string) (uid string) {
	defer aria.resp.Reset()
	msg := `{
  "id": "QXJpYU5nXzE2NDA2MzAyMjNfMC43MzkwOTc4NTA5NTYyMDQx",
  "jsonrpc": "2.0",
  "method": "aria2.addUri",
  "params": [
    [
      "%s"
    ],
	{
	"dir":"%s",
	"out":"%s",
	"referer":"%s"
	}
  ]
}`
	msg = fmt.Sprintf(msg, url, dir, out, referer)

	aria.req.SetBodyString(msg)

	fasthttp.Do(aria.req, aria.resp)

	parse := gjson.ParseBytes(aria.resp.Body())
	return parse.Get("result").String()
}

func (aria *AriaEngine) getStatus(uid string) (status string) {
	defer aria.resp.Reset()

	msg := `{"jsonrpc":"2.0","method":"aria2.tellStatus","id":"QXJpYU5nXzE2NDA2MzMxNzhfMC4wNDgxOTE2OTU0MjQ0MDkwNA==","params":["%s"]}`
	msg = fmt.Sprintf(msg, uid)

	aria.req.SetBodyString(msg)

	fasthttp.Do(aria.req, aria.resp)

	parse := gjson.ParseBytes(aria.resp.Body())
	return parse.Get("result.status").String()
}

func (aria *AriaEngine) IsDone(uid string) bool {
	return aria.getStatus(uid) == "complete"
}

type GlobalStat struct {
	DownloadSpeed   int64 `json:"downloadSpeed"`
	NumActive       int64  `json:"numActive"`
	NumStopped      int64 `json:"numStopped"`
	NumStoppedTotal int64 `json:"numStoppedTotal"`
	NumWaiting      int64 `json:"numWaiting"`
	UploadSpeed     int64 `json:"uploadSpeed"`
}

func (aria *AriaEngine) GetGlobalStat() GlobalStat {
	defer aria.resp.Reset()
	msg := `{"jsonrpc":"2.0","method":"aria2.getGlobalStat","id":"QXJpYU5nXzE2NDA2MzMxNzhfMC4wNDgxOTE2OTU0MjQ0MDkwNA=="}`

	aria.req.SetBodyString(msg)

	fasthttp.Do(aria.req, aria.resp)

	parse := gjson.ParseBytes(aria.resp.Body())
	result := parse.Get("result")
	return GlobalStat{
		DownloadSpeed:   result.Get("downloadSpeed").Int(),
		NumActive:       result.Get("numActive").Int(),
		NumStopped:      result.Get("numStopped").Int(),
		NumStoppedTotal: result.Get("numStoppedTotal").Int(),
		NumWaiting:      result.Get("numWaiting").Int(),
		UploadSpeed:     result.Get("uploadSpeed").Int(),
	}
}
