package main

import (
	"awesomeProject/config"
	"awesomeProject/server"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path"
)

// var forceFlagValue bool
var directionFlagValue string

func init() {

	//	flag.BoolVar(&forceFlagValue, "f", false, "Optional. Force overwrite if given error about save.")
	flag.StringVar(&directionFlagValue, "d", "", "Required. Determines if a save will be uploaded (up) or downloaded (down)")
}

func main() {
	flag.Parse()
	c, err := config.ReadConfig("./user.config")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	switch directionFlagValue {

	case "up":
		filepath := path.Join(savePath, fmt.Sprintf("ER0000%s", extension))
		upload, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = server.UploadSave(c.Name, c.Key, "ER0000", upload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "down":

		filepath := path.Join(savePath, fmt.Sprintf("ER0000%s", extension))
		dl, err := server.DownloadSave(c.Name, c.Key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b, err := base64.StdEncoding.DecodeString(dl)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile(filepath, b, 0644)

	default:
		flag.PrintDefaults()
	}

}
