//go:build cmd
// +build cmd

package main

import (
	"fmt"
	"m3u8/initial"
)

func main() {
	initial.InitLogger("dev")

	var url string
	for {
		fmt.Printf("url :")
		if _, err := fmt.Scan(&url); err != nil {
			fmt.Println(err)
			continue
		}
		go initial.Down(url)
	}
}
