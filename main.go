package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
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

var a = make(chan int, 1)

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateRunningTask
)

var appState ApplicationState = ApplicationStateIdle

func init() {

	flag.StringVar(&directionFlagValue, "d", "", "Skip GUI and determines if a save will be uploaded (up) or downloaded (down)")
	flag.BoolVar(&installFlagValue, "i", false, "Skip GUI and install latest Seamless Coop.")
}

func toggleIdle() {
	if appState == ApplicationStateIdle {
		appState = ApplicationStateRunningTask
		a <- 1
	} else {
		appState = ApplicationStateIdle
		a <- 2
	}
}

func saveCommand() {
	toggleIdle()
	defer toggleIdle()
	c := config.NewConfig(configPath)

	c.Name = widgets.GetName()
	c.Key = widgets.GetKey()
	c.FileExtension = widgets.GetExt()
	c.Path = widgets.GetSavePath()
	c.GamePath = widgets.GetGamePath()
	err := c.Save()
	if err != nil {
		log.Fatal(err)
	}
}

func reloadCommand() {

	toggleIdle()
	defer toggleIdle()
	c, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	widgets.SetConfig(c)
}

func uploadCommand() {

	toggleIdle()
	defer toggleIdle()
	executeUpload()
}

func downloadCommand() {

	toggleIdle()
	defer toggleIdle()
	executeDownload()
}

func installCommand() {
	toggleIdle()
	defer toggleIdle()
	executeInstall()

}

func getSystemSavePath(c config.Config) string {
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
	go reloadCommand()
	for {
		switch e := window.Event().(type) {
		case transfer.DataEvent:
			log.Println("GOT DATA EVENT")
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:

			gtx := app.NewContext(&ops, e)

			// invalidate the window when the channel has been written to. Allows the app to switch between the 2 screens.
			go func() {
				select {
				case <-a:
					window.Invalidate()
				}
			}()

			if appState == ApplicationStateIdle {
				widgets.StandardLayout(theme, gtx)
			} else {
				layout.Stack{Alignment: layout.Center}.Layout(gtx,
					layout.Expanded(func(g layout.Context) layout.Dimensions {

						h := material.H3(theme, "Running...")
						h.Alignment = text.Middle
						return h.Layout(g)
					}))
			}
			e.Frame(gtx.Ops)
		}
	}
}

func executeInstall() {
	c, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = server.InstallCoop(c.GamePath)
	if err != nil {
		log.Fatal(err)
	}
}
func executeUpload() {
	c, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	filepath := getSystemSavePath(c)
	upload, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	err = server.UploadSave(c.Name, c.Key, "ER0000", upload)
	if err != nil {
		log.Fatal(err)
	}
}

func executeDownload() {

	c, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	filepath := getSystemSavePath(c)
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
}
func main() {
	flag.Parse()

	// Wire up the widgets buttons
	widgets.SaveCommand = saveCommand
	widgets.ReloadCommand = reloadCommand
	widgets.UploadCommand = uploadCommand
	widgets.DownloadCommand = downloadCommand
	widgets.InstallCommand = installCommand

	switch directionFlagValue {
	case "up": //Passed -d up to application
		executeUpload()
	case "down": //Passed -d down to application
		executeDownload()
	default:
		go func() {
			window := new(app.Window)
			window.Option(app.Title("Lost Grace"))
			window.Option(app.Size(unit.Dp(1234), unit.Dp(950)))
			err := run(window)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		app.Main()
	}
}
