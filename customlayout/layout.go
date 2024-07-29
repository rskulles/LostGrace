package customlayout

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"lostgrace/config"
)

// The Widgets
var (
	nameInput     = &widget.Editor{SingleLine: true, Submit: true}
	keyInput      = &widget.Editor{SingleLine: true, Submit: true}
	savePathInput = &widget.Editor{SingleLine: true, Submit: true}
	extInput      = &widget.Editor{SingleLine: true, Submit: true}
	gamePathInput = &widget.Editor{SingleLine: true, Submit: true}

	saveConfigButton   = new(widget.Clickable)
	reloadConfigButton = new(widget.Clickable)
	uploadSaveButton   = new(widget.Clickable)
	downloadSaveButton = new(widget.Clickable)
	installButton      = new(widget.Clickable)
)

// Button Click Handlers
var (
	SaveCommand     func() = nil
	ReloadCommand   func() = nil
	UploadCommand   func() = nil
	DownloadCommand func() = nil
	InstallCommand  func() = nil
)

func buttonHandleCommand(gtx layout.Context, w *widget.Clickable, f func()) {
	if w.Clicked(gtx) && f != nil {
		go f()
	}
}

func SetName(name string) {
	nameInput.SetText(name)
}

func SetKey(key string) {
	keyInput.SetText(key)
}

func SetSavePath(path string) {
	savePathInput.SetText(path)
}

func SetExt(ext string) {
	extInput.SetText(ext)
}

func SetGamePath(gamePath string) {
	gamePathInput.SetText(gamePath)
}

func GetName() string {
	return nameInput.Text()
}

func GetSavePath() string {
	return savePathInput.Text()
}
func GetKey() string {
	return keyInput.Text()
}

func GetExt() string {
	return extInput.Text()
}

func GetGamePath() string {
	return gamePathInput.Text()
}

func SetConfig(c config.Config) {

	SetName(c.Name)
	SetKey(c.Key)
	SetSavePath(c.Path)
	SetExt(c.FileExtension)
	SetGamePath(c.GamePath)
}
func StandardLayout(theme *material.Theme, gtx layout.Context) {

	buttonHandleCommand(gtx, saveConfigButton, SaveCommand)
	buttonHandleCommand(gtx, reloadConfigButton, ReloadCommand)
	buttonHandleCommand(gtx, uploadSaveButton, UploadCommand)
	buttonHandleCommand(gtx, downloadSaveButton, DownloadCommand)
	buttonHandleCommand(gtx, installButton, InstallCommand)

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
				header := material.H4(theme, "Config")
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
				e := material.Editor(theme, savePathInput, "Set Save Path")
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
			layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				header := material.Label(theme, unit.Sp(25), "Elden Ring \"Game\" Path")
				return header.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(theme, gamePathInput, "Set Game Path. Right click Elden Ring in your Steam library. Manage->Browse Local Files. Copy/Paste/Type the path into this box.")
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
				h := material.H4(theme, "Sync")
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
				h := material.H4(theme, "Install Latest ERSC")
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

					btn := material.Button(theme, installButton, "Install")
					btn.Inset.Top = unit.Dp(25)
					btn.Inset.Bottom = unit.Dp(25)
					return btn.Layout(gtx)
				})
			}),
		)
	})
}
