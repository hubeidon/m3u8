//go:build gui
//+build gui

package main

import (
	"gioui.org/app"
	"m3u8/initial"
)

func main() {
	initial.InitLogger("dev")

	// go gui()

	app.Main()
}
