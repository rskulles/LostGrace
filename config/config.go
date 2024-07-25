package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Name          string
	Key           string
	Path          string
	FileExtension string
}

func ReadConfig(path string) (Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		//Check if comment
		if strings.HasPrefix(text, "#") {
			continue
		}
		split := strings.Split(text, "=")
		switch strings.TrimSpace(split[0]) {
		case "name":
			config.Name = strings.TrimSpace(split[1])
		case "key":
			config.Key = strings.TrimSpace(split[1])
		case "path":
			config.Path = strings.TrimSpace(split[1])
		case "extension":
			config.FileExtension = strings.TrimSpace(split[1])
		default:
			continue
		}
	}

	return config, nil
}
