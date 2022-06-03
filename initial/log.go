package initial

import (
	"io"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func InitLogger(level string) {
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	f, _ := os.OpenFile("_m3u8.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)

	if level == "dev"{
		log.SetReportCaller(true)
		log.SetOutput(io.MultiWriter(f, os.Stdout))
	}else{
		log.SetOutput(f)
	}
}
