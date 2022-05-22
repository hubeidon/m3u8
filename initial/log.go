package initial

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//log.SetReportCaller(true)

	f, _ := os.OpenFile("_m3u8.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)

	//log.SetOutput(io.MultiWriter(f, os.Stdout))
	//
	//log.Println("test")

	log.SetOutput(f)
}
