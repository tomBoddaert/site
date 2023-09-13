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

func buildTemplated(config *Config, dst_mode fs.FileMode) {
	dir, err := os.ReadDir(config.TemplatedSrcDir)
	check(err)

	for _, src_dir_meta := range dir {
		if !src_dir_meta.IsDir() {
			logger.SetOutput(os.Stderr)
			logger.Warnf("%v is not a directory! Ignoring", path.Join(config.TemplatedSrcDir, src_dir_meta.Name()))
			logger.SetOutput(os.Stdout)
			continue
		}

		tmpl_name := path.Join(config.TemplateDir, src_dir_meta.Name())

		logger.Debugf("Building template (%v)", tmpl_name)
		tmpl := getTemplate(tmpl_name)

		config := BuildConfig{
			Config:       config,
			DstMode:      dst_mode,
			Template:     tmpl,
			Data:         getData(config.DataFile),
			TemplateName: src_dir_meta.Name(),
			SubDir:       "",
		}

		src_dir_path := path.Join(path.Join(config.Config.TemplatedSrcDir, config.TemplateName))
		src_dir, err := os.ReadDir(src_dir_path)
		check(err)

		logger.Debugf("Building directory (%v -> %v)", src_dir_path, config.Config.DstDir)

		for _, sub_src := range src_dir {
			config.Entry = sub_src
			buildTemplatedRecursive(config)
		}
	}
}

func buildTemplatedRecursive(config BuildConfig) {
	src_path := config.SrcPath()
	dst_path := config.DstPath()

	if config.Entry.IsDir() {
		logger.Debugf("Building directory (%v -> %v)", src_path, dst_path)

		err := os.Mkdir(dst_path, config.DstMode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			check(err)
		}

		dir, err := os.ReadDir(src_path)
		check(err)

		config.AppendSubSrc(config.Entry.Name())

		for _, sub := range dir {
			config.Entry = sub
			buildTemplatedRecursive(config)
		}
	} else {
		logger.Debugf("Building file (%v -> %v)", src_path, dst_path)

		content, err := os.ReadFile(src_path)
		check(err)

		content_tmpl, err := createTemplate("Content").Parse(string(content))
		check(err)

		new_tmpl, err := config.Template.Clone()
		check(err)

		new_tmpl.AddParseTree("Content", content_tmpl.Tree)

		dst, err := os.Create(dst_path)
		check(err)
		defer dst.Close()

		err = dst.Chmod(config.DstMode)
		check(err)

		data := buildData(&config)

		if config.Config.FmtTemplatedHtml && path.Ext(config.Entry.Name()) == ".html" {
			logger.Debug("Formatting HTML")

			preformat := new(bytes.Buffer)
			err := new_tmpl.Execute(preformat, data)
			if err != nil {
				logger.Errorf("Error executing template (%v) on file (%v)", config.TemplateName, src_path)
				logger.Infof("Error: %v", err)
			}

			formatted := gohtml.FormatBytes(preformat.Bytes())
			dst.Write(formatted)
		} else {
			err := new_tmpl.Execute(dst, data)
			if err != nil {
				logger.Errorf("Error executing template (%v) on file (%v)", config.TemplateName, src_path)
				logger.Infof("Error: %v", err)
			}
		}
	}
}
