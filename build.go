package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/yosssi/gohtml"
)

type BuildConfig struct {
	Config       *Config
	DstMode      fs.FileMode
	Template     *template.Template
	Data         map[string]map[string]any
	TemplateName string
	SubDir       string
	Entry        os.DirEntry
}

func (config *BuildConfig) SrcPath() string {
	return path.Join(config.Config.TemplatedSrcDir, config.TemplateName, config.SubDir, config.Entry.Name())
}

func (config *BuildConfig) DstPath() string {
	return path.Join(config.Config.DstDir, config.SubDir, config.Entry.Name())
}

func (config *BuildConfig) DataPath() string {
	return path.Join(config.TemplateName, config.SubDir, config.Entry.Name())
}

func (config *BuildConfig) AppendSubSrc(frag string) {
	config.SubDir = path.Join(config.SubDir, frag)
}

func buildTemplated(config *Config, dstMode fs.FileMode) {
	dir, err := os.ReadDir(config.TemplatedSrcDir)
	check(err)

	for _, srcDirMeta := range dir {
		if !srcDirMeta.IsDir() {
			logger.SetOutput(os.Stderr)
			logger.Warnf("%v is not a directory! Ignoring", path.Join(config.TemplatedSrcDir, srcDirMeta.Name()))
			logger.SetOutput(os.Stdout)
			continue
		}

		tmplName := path.Join(config.TemplateDir, srcDirMeta.Name())

		logger.Debugf("Building template (%v)", tmplName)
		tmpl := getTemplate(tmplName)

		config := BuildConfig{
			Config:       config,
			DstMode:      dstMode,
			Template:     tmpl,
			Data:         getData(config.DataFile),
			TemplateName: srcDirMeta.Name(),
			SubDir:       "",
		}

		srcDirPath := path.Join(path.Join(config.Config.TemplatedSrcDir, config.TemplateName))
		srcDir, err := os.ReadDir(srcDirPath)
		check(err)

		logger.Debugf("Building directory (%v -> %v)", srcDirPath, config.Config.DstDir)

		for _, subSrc := range srcDir {
			config.Entry = subSrc
			buildTemplatedRecursive(config)
		}
	}
}

func buildTemplatedRecursive(config BuildConfig) {
	srcPath := config.SrcPath()
	dstPath := config.DstPath()

	if config.Entry.IsDir() {
		logger.Debugf("Building directory (%v -> %v)", srcPath, dstPath)

		err := os.Mkdir(dstPath, config.DstMode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			check(err)
		}

		dir, err := os.ReadDir(srcPath)
		check(err)

		config.AppendSubSrc(config.Entry.Name())

		for _, sub := range dir {
			config.Entry = sub
			buildTemplatedRecursive(config)
		}
	} else {
		logger.Debugf("Building file (%v -> %v)", srcPath, dstPath)

		content, err := os.ReadFile(srcPath)
		check(err)

		contentTmpl, err := createTemplate("Content").Parse(string(content))
		check(err)

		newTmpl, err := config.Template.Clone()
		check(err)

		newTmpl.AddParseTree("Content", contentTmpl.Tree)

		dst, err := os.Create(dstPath)
		check(err)
		defer dst.Close()

		err = dst.Chmod(config.DstMode)
		check(err)

		data := buildData(&config)

		if config.Config.FmtTemplatedHTML && path.Ext(config.Entry.Name()) == ".html" {
			logger.Debug("Formatting HTML")

			preformat := new(bytes.Buffer)
			err := newTmpl.Execute(preformat, data)
			if err != nil {
				logger.Errorf("Error executing template (%v) on file (%v)", config.TemplateName, srcPath)
				logger.Infof("Error: %v", err)
			}

			formatted := gohtml.FormatBytes(preformat.Bytes())
			dst.Write(formatted)
		} else {
			err := newTmpl.Execute(dst, data)
			if err != nil {
				logger.Errorf("Error executing template (%v) on file (%v)", config.TemplateName, srcPath)
				logger.Infof("Error: %v", err)
			}
		}
	}
}
