package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"

	"github.com/yosssi/gohtml"
)

type RawConfig struct {
	Config  *Config
	DstMode os.FileMode
	SubSrc  string
	Entry   os.DirEntry
}

func (config *RawConfig) SrcPath() string {
	return path.Join(config.Config.RawSrcDir, config.SubSrc, config.Entry.Name())
}

func (config *RawConfig) DstPath() string {
	return path.Join(config.Config.DstDir, config.SubSrc, config.Entry.Name())
}

func (config *RawConfig) AppendSubSrc(frag string) {
	config.SubSrc = path.Join(config.SubSrc, frag)
}

func copyRaw(config *Config, dstMode os.FileMode) {
	dir, err := os.ReadDir(config.RawSrcDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		logger.Warnf("No raw directory found (%v)", config.RawSrcDir)
		return
	}
	check(err)

	for _, sub := range dir {
		copyRawRecursive(RawConfig{
			Config:  config,
			DstMode: dstMode,
			SubSrc:  "",
			Entry:   sub,
		})
	}
}

func copyRawRecursive(config RawConfig) {
	srcPath := config.SrcPath()
	dstPath := config.DstPath()

	if config.Entry.IsDir() {
		logger.Debugf("Copying directory (%v -> %v)", srcPath, dstPath)

		err := os.Mkdir(dstPath, config.DstMode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			check(err)
		}

		dir, err := os.ReadDir(srcPath)
		check(err)

		config.AppendSubSrc(config.Entry.Name())

		for _, sub := range dir {
			config.Entry = sub
			copyRawRecursive(config)
		}
	} else {
		if !config.Config.IncludeTS && path.Ext(config.Entry.Name()) == ".ts" {
			logger.Debugf("Ignoring TS file (%v)", srcPath)
			return
		}

		logger.Debugf("Copying file (%v -> %v)", srcPath, dstPath)

		src, err := os.Open(srcPath)
		check(err)
		defer src.Close()

		dst, err := os.Create(dstPath)
		check(err)
		defer dst.Close()

		err = dst.Chmod(config.DstMode)
		check(err)

		if config.Config.FmtRawHTML && path.Ext(config.Entry.Name()) == ".html" {
			logger.Debug("Formatting HTML file")

			preformat := new(bytes.Buffer)
			_, err := io.Copy(preformat, src)
			check(err)

			formatted := gohtml.FormatBytes(preformat.Bytes())
			dst.Write(formatted)
		} else {
			_, err = io.Copy(dst, src)
			check(err)
		}
	}
}
