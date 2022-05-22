package main

import (
	"fmt"
	"m3u8/initial"
	"os"
	"time"

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
		app.Size(unit.Dp(400), unit.Dp(230)),
		// app.MaxSize(unit.Dp(400), unit.Dp(170)),
		app.MinSize(unit.Dp(400), unit.Dp(230)),
	)
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	var Input widget.Editor

	var speed int

	go func() {
		for {
			speed++
			time.Sleep(time.Second)
		}
	}()

	// th defnes the material design style
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
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bar := material.Label(th, unit.Dp(25), fmt.Sprintf("speed : %d",speed))
							return bar.Layout(gtx)
						})
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if startButton.Clicked() {
							initial.HttpOrLocal(Input.Text())
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

	//close window
	initial.CloseWindow()
	os.Exit(0)
}
