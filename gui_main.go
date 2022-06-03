//go:build gui
// +build gui

package main

import (
	"gioui.org/app"
	"m3u8/initial"
)

func main() {
	initial.Run("")

	go gui()

	app.Main()
}
