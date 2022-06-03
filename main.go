//go:build cmd
//+build cmd

package main

import (
	"fmt"
	"m3u8/initial"
	"runtime"
	"strings"
	"time"
)

func PrintGONum() {
	for {
		fmt.Println(runtime.NumGoroutine())
		time.Sleep(time.Second * 5)
	}
}

func main() {
	initial.Run("")

	go PrintGONum()

	var url string
	for {
		fmt.Printf("url :")
		if _, err := fmt.Scan(&url); err != nil {
			fmt.Println(err)
			continue
		}
		go initial.HttpOrLocal(strings.TrimSpace(url))
	}
}
