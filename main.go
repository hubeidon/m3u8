package main

import (
	"m3u8/initial"

	"gioui.org/app"
)

func main() {
	go gui()

	initial.Run()

	app.Main()
}