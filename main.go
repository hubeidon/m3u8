package main

import (
	"m3u8/initial"

	"gioui.org/app"
)

func main() {
	initial.Run("")

	go gui()

	app.Main()
}
