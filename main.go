package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"log"
	"lostgrace/config"
	"lostgrace/server"
	"lostgrace/widgets"
	"os"
	"path"
)

var directionFlagValue string
var installFlagValue bool

const configPath = "./user.config"

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateRunningTask
)

func init() {

	flag.StringVar(&directionFlagValue, "d", "", "Skip GUI and determines if a save will be uploaded (up) or downloaded (down)")
	flag.BoolVar(&installFlagValue, "i", false, "Skip GUI and install latest Seamless Coop.")
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
	appState := ApplicationStateIdle
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if appState == ApplicationStateIdle {
				widgets.StandardLayout(theme, gtx)
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
			window.Option(app.Size(unit.Dp(600), unit.Dp(950)))
			err := run(window)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		app.Main()
	}

}
