//go:build cmd
// +build cmd

package main

import (
	"fmt"
	"m3u8/initial"
)

// var u string = `https://webplay.weilekangnet.com:59666/data6/B1AF27ADA87782E7/E6387F593401CC76/play.ts?_KS=24b7da2d31ebde13ea5e0a530875cb1b&_KE=1654904339`

func main() {
	initial.InitLogger("dev")

	var url string
	for {
		fmt.Printf("url :")
		if _, err := fmt.Scan(&url); err != nil {
			fmt.Println(err)
			continue
		}
		// go initial.HttpOrLocal(strings.TrimSpace(url))
	}
}
