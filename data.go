package main

import (
	"encoding/json"
	"errors"
	"os"
)

func getData(path string) map[string]map[string]any {
	data := make(map[string]map[string]any)

	file, err := os.Open(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			check(err)
		}

		return data
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	check(err)

	return data
}

func buildData(config *BuildConfig) map[string]any {
	data := make(map[string]any)

	for key, value := range config.Data["default"] {
		data[key] = value
	}

	for key, value := range config.Data[config.TemplateName] {
		data[key] = value
	}

	for key, value := range config.Data[config.DataPath()] {
		data[key] = value
	}

	return data
}
