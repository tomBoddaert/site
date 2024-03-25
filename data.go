package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

func GetData(path string) map[string]map[string]any {
	data := make(map[string]map[string]any)

	file, err := os.Open(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			check(err)
		}

		logger.Debugf("No data found (%v)", path)

		return data
	}
	defer file.Close()

	logger.Debugf("Found data (%v)", path)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	check(err)

	return data
}

func BuildData(config *buildConfig, path_ string) map[string]any {
	data := make(map[string]any)

	for key, value := range config.Data["default"] {
		data[key] = value
	}

	for key, value := range config.Data[config.TemplateName] {
		data[key] = value
	}

	dataPath := path.Join(config.TemplateName, path_)
	for key, value := range config.Data[dataPath] {
		data[key] = value
	}

	return data
}
