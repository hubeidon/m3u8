// go:build gui
// +build gui

package main

import (
	"fmt"
	"m3u8/initial"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func gui() {
	w := app.NewWindow(
		app.Title("m3u8"),
		app.Size(unit.Dp(400), unit.Dp(300)),
		// app.MaxSize(unit.Dp(400), unit.Dp(170)),
		app.MinSize(unit.Dp(400), unit.Dp(300)),
	)

	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	var Input widget.Editor

	var (
		downloadSpeed, numActive, numWaiting int64
	)

	// th defnes the material design style
	//th := material.NewTheme(gofont.Collection())
	th := material.NewTheme(gofont.Collection())

	for e := range w.Events() {

		switch e := e.(type) {
		case system.FrameEvent:

			gtx := layout.NewContext(&ops, e)
			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Editor(th, &Input, "url") // Here progress is used
							return bar.Layout(gtx)
						})
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							// Top:    unit.Dp(25),
							// Bottom: unit.Dp(25),
							Right: unit.Dp(35),
							Left:  unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Label(th, unit.Dp(20), fmt.Sprintf("speed : %.2fM/s", float64(downloadSpeed/1024/1024)))
							return bar.Layout(gtx)
						})
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							// Top:    unit.Dp(25),
							// Bottom: unit.Dp(25),
							Right: unit.Dp(35),
							Left:  unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Label(th, unit.Dp(20), fmt.Sprintf("active : %d", numActive))
							return bar.Layout(gtx)
						})
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							// Top:    unit.Dp(25),
							// Bottom: unit.Dp(25),
							Right: unit.Dp(35),
							Left:  unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Label(th, unit.Dp(20), fmt.Sprintf("wait : %d", numWaiting))
							return bar.Layout(gtx)
						})
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if startButton.Clicked() {
							text := Input.Text()
							if text != "" {
								go initial.Down(text)
							}
							Input.SetText("")
						}

						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Button(th, &startButton, "start") // Here progress is used
							return bar.Layout(gtx)
						})
					},
				),
			)
			e.Frame(&ops)
		}
	}

	os.Exit(0)
}
