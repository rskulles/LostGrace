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
	GamePath      string
	configPath    string
}

func ReadConfig(path string) (Config, error) {
	config := Config{configPath: path}
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
		case "save_path":
			config.Path = strings.TrimSpace(split[1])
		case "extension":
			config.FileExtension = strings.TrimSpace(split[1])
		case "game_path":
			config.GamePath = strings.TrimSpace(split[1])
		default:
			continue
		}
	}

	return config, nil
}

func NewConfig(path string) Config {
	return Config{configPath: path}
}

func (config *Config) Save() error {
	file, err := os.Create(config.configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	_, _ = fileWriter.WriteString(fmt.Sprintf("name=%s\n", config.Name))
	_, _ = fileWriter.WriteString(fmt.Sprintf("key=%s\n", config.Key))
	_, _ = fileWriter.WriteString(fmt.Sprintf("save_path=%s\n", config.Path))
	_, _ = fileWriter.WriteString(fmt.Sprintf("extension=%s\n", config.FileExtension))
	_, _ = fileWriter.WriteString(fmt.Sprintf("game_path=%s\n", config.GamePath))
	_ = fileWriter.Flush()
	return nil
}
