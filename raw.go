package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/yosssi/gohtml"
)

func CopyRaw(config *Config) {
	root := os.DirFS(config.RawSrcDir)
	err := fs.WalkDir(root, ".", config.copyRawWalk)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		logger.Warnf("No raw directory found (%v)", config.RawSrcDir)
		return
	}
}

func (config *Config) copyRawWalk(path_ string, d fs.DirEntry, err error) error {
	excludeMatch, excludeIndex := PathMatch(path_, config.ExcludePaths)
	if excludeMatch != nil {
		logger.Debugf("Excluding path (%v): rule %v", *excludeMatch, excludeIndex+1)

		if d.IsDir() {
			return fs.SkipDir
		} else {
			return nil
		}
	}

	check(err)

	srcPath := path.Join(config.RawSrcDir, path_)
	dstPath := path.Join(config.DstDir, path_)

	if d.IsDir() {
		logger.Debugf("Copying directory (%v -> %v)", srcPath, dstPath)

		err := os.Mkdir(dstPath, config.DstMode.FileMode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			check(err)
		}
	} else {
		logger.Debugf("Copying file (%v -> %v)", srcPath, dstPath)

		src, err := os.Open(srcPath)
		check(err)
		defer src.Close()

		dst, err := os.Create(dstPath)
		check(err)
		defer dst.Close()

		err = dst.Chmod(config.DstMode.FileMode)
		check(err)

		if config.FmtRawHTML && path.Ext(d.Name()) == ".html" {
			logger.Debug("Formatting HTML file")

			preformat := new(bytes.Buffer)
			_, err := io.Copy(preformat, src)
			check(err)

			formatted := gohtml.FormatBytes(preformat.Bytes())
			_, err = dst.Write(formatted)
			check(err)
		} else {
			_, err := io.Copy(dst, src)
			check(err)
		}
	}

	return nil
}
