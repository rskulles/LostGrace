package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

var (
	nameInput = &widget.Editor{SingleLine: true, Submit: true}
	keyInput  = &widget.Editor{SingleLine: true, Submit: true}
	pathInput = &widget.Editor{SingleLine: true, Submit: true}
	extInput  = &widget.Editor{SingleLine: true, Submit: true}

	gamePathInput      = &widget.Editor{SingleLine: true, Submit: true}
	saveConfigButton   = new(widget.Clickable)
	reloadConfigButton = new(widget.Clickable)
	uploadSaveButton   = new(widget.Clickable)
	downloadSaveButton = new(widget.Clickable)
	installButton      = new(widget.Clickable)
)

func StandardLayout(theme *material.Theme, gtx layout.Context) {
	appMargins := layout.Inset{
		Top:    unit.Dp(0),
		Bottom: unit.Dp(0),
		Left:   unit.Dp(10),
		Right:  unit.Dp(10),
	}
	appMargins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceStart,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.H3(theme, "Config")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Name")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, nameInput, "Set Name")
				b := widget.Border{Color: color.NRGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
				return b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Key")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, keyInput, "Set Key")
				b := widget.Border{Color: color.NRGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
				return b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Save Path")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, pathInput, "Set Save Path")
				b := widget.Border{Color: color.NRGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
				return b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Extension")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, extInput, "Set Extension")
				b := widget.Border{Color: color.NRGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
				return b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {

				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}

				return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					btn := material.Button(theme, saveConfigButton, "Save")
					btn2 := material.Button(theme, reloadConfigButton, "Reload")
					return SplitVisual{}.Layout(gtx, unit.Dp(50), btn.Layout, btn2.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				h := material.H3(theme, "Sync")
				return h.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {

				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}

				return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					btn := material.Button(theme, uploadSaveButton, "Upload")
					btn2 := material.Button(theme, downloadSaveButton, "Download")
					return SplitVisual{}.Layout(gtx, unit.Dp(50), btn.Layout, btn2.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				h := material.H3(theme, "Install Latest ERSC")
				return h.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Save Path")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, gamePathInput, "Set Game Path")
				b := widget.Border{Color: color.NRGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
				return b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}
				return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(theme, installButton, "Install")
					return btn.Layout(gtx)

				})
			}),
		)
	})
}
