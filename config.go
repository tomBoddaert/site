package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type Config struct {
	TemplateDir      string
	TemplateRootFile string
	TemplatedSrcDir  string
	RawSrcDir        string
	DstDir           string
	DstMode          FileMode
	DataFile         string
	FmtTemplatedHTML bool `json:"FmtTemplatedHtml"`
	FmtRawHTML       bool `json:"FmtRawHtml"`
	ExcludePaths     []Regexp
	PrebuildCmds     [][]string
	PostbuildCmds    [][]string

	NotFoundPath string
}

const configPath = "siteConfig.json"

var defaultConfig = Config{
	TemplateDir:      "templates",
	TemplateRootFile: "root.html",
	TemplatedSrcDir:  "pages",
	RawSrcDir:        "raw",
	DstDir:           "out",
	DstMode:          FileMode{fs.FileMode(0755)},
	DataFile:         "data.json",
	FmtTemplatedHTML: false,
	FmtRawHTML:       false,
	ExcludePaths:     []Regexp{},
	PrebuildCmds:     [][]string{},
	PostbuildCmds:    [][]string{},

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
