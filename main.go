package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"m3u8/initial"
	"os"
	"os/signal"
)

func main() {
	initial.Run("")

	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go listen(c)

	var url string
	for {
		fmt.Printf("url :")
		if _, err := fmt.Scan(&url); err != nil {
			fmt.Println(err)
			continue
		}
		go initial.Down(url)
		log.Infof("start downlaod %s", url[62:78])
	}
}

func listen(c chan os.Signal) {
	for {
		<-c
		if err := os.RemoveAll("./data"); err != nil {
			log.Error(err)
		} else {
			fmt.Println("\n已删除缓存数据, 退出程序!")
			os.Exit(0)
		}
	}
}
