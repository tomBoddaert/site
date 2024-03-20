package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	TemplateDir      string
	TemplatedSrcDir  string
	RawSrcDir        string
	DstDir           string
	DstMode          string
	DataFile         string
	FmtTemplatedHTML bool `json:"FmtTemplatedHtml"`
	FmtRawHTML       bool `json:"FmtRawHtml"`
	TranspileTS      bool
	TSArgs           []string
	IncludeTS        bool

	NotFoundPath string
}

const configPath = "siteConfig.json"

var defaultConfig = Config{
	TemplateDir:      "templates",
	TemplatedSrcDir:  "pages",
	RawSrcDir:        "raw",
	DstDir:           "out",
	DstMode:          "0755",
	DataFile:         "data.json",
	FmtTemplatedHTML: false,
	FmtRawHTML:       false,
	TranspileTS:      true,
	TSArgs:           []string{},
	IncludeTS:        false,

	NotFoundPath: "404",
}

func getConfig() Config {
	file, err := os.Open(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			check(err)
		}

		logger.Debugf("No config found (%v)", configPath)

		return defaultConfig
	}

	defer file.Close()

	logger.Debugf("Found config (%v)", configPath)

	decoder := json.NewDecoder(file)

	config := defaultConfig
	err = decoder.Decode(&config)
	check(err)

	return config
}
