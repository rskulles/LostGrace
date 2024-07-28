package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"lostgrace/config"
	"lostgrace/server"
	"lostgrace/widgets"
	"os"
	"path"
)

// var forceFlagValue bool
var directionFlagValue string

const configPath = "./user.config"

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateRunningTask
)

func init() {

	//	flag.BoolVar(&forceFlagValue, "f", false, "Optional. Force overwrite if given error about save.")
	flag.StringVar(&directionFlagValue, "d", "", "Required. Determines if a save will be uploaded (up) or downloaded (down)")
}

func getSavePath(c config.Config) string {
	extension := ""
	if c.FileExtension == "default" {
		extension = ".co2"
	} else {
		extension = c.FileExtension
	}

	savePath := ""
	if c.Path == "default" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		savePath = path.Join(userConfigDir, "EldenRing")
		entries, err := os.ReadDir(savePath)
		directoryCount := 0
		directoryName := ""
		for _, entry := range entries {
			if entry.IsDir() {
				directoryCount++
				directoryName = entry.Name()
			}
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if directoryCount == 1 {
			savePath = path.Join(savePath, directoryName)
		} else {
			fmt.Println("Multiple directories found under your save directory. Please set path variable directly in user.config")
		}
	} else {
		savePath = c.Path
	}
	filepath := path.Join(savePath, fmt.Sprintf("ER0000%s", extension))
	return filepath
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	nameInput := &widget.Editor{SingleLine: true, Submit: true}
	keyInput := &widget.Editor{SingleLine: true, Submit: true}
	pathInput := &widget.Editor{SingleLine: true, Submit: true}
	extInput := &widget.Editor{SingleLine: true, Submit: true}

	saveConfigButton := new(widget.Clickable)
	reloadConfigButton := new(widget.Clickable)
	uploadSaveButton := new(widget.Clickable)
	downloadSaveButton := new(widget.Clickable)
	appState := ApplicationStateIdle
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if appState == ApplicationStateIdle {

				layout.Flex{
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
							return widgets.SplitVisual{}.Layout(gtx, unit.Dp(50), btn.Layout, btn2.Layout)
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
							return widgets.SplitVisual{}.Layout(gtx, unit.Dp(50), btn.Layout, btn2.Layout)
						})
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
				)
			}
			e.Frame(gtx.Ops)
		}
	}
}

func main() {
	flag.Parse()
	c, err := config.ReadConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}

	switch directionFlagValue {
	case "up": //Passed -d up to application
		filepath := getSavePath(c)
		upload, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		err = server.UploadSave(c.Name, c.Key, "ER0000", upload)
		if err != nil {
			log.Fatal(err)
		}
	case "down": //Passed -d down to application
		filepath := getSavePath(c)
		dl, err := server.DownloadSave(c.Name, c.Key)

		if err != nil {
			log.Fatal(err)
		}

		//TODO: This works, but do it as part of retrieval and return the bytes.
		b, err := base64.StdEncoding.DecodeString(dl)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filepath, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	default:
		go func() {
			window := new(app.Window)
			window.Option(app.Title("Lost Grace"))
			window.Option(app.Size(unit.Dp(600), unit.Dp(800)))
			err := run(window)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		app.Main()
	}

}
