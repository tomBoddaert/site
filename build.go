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

type buildConfig struct {
	*Config
	*template.Template
	Data         map[string]map[string]any
	TemplateName string
}

func BuildTemplated(config *Config) {
	dir, err := os.ReadDir(config.TemplatedSrcDir)
	check(err)

	data := GetData(config.DataFile)

	for _, srcDirMeta := range dir {
		if !srcDirMeta.IsDir() {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Template content is not a directory (%v)", path.Join(config.TemplatedSrcDir, srcDirMeta.Name()))
			logger.SetOutput(os.Stdout)

			continue
		}

		tmplPath := path.Join(config.TemplateDir, srcDirMeta.Name())

		logger.Debugf("Building template (%v)", tmplPath)
		tmpl := GetTemplate(config, srcDirMeta.Name())
		if tmpl == nil {
			continue
		}

		config := buildConfig{
			Config:       config,
			Template:     tmpl,
			Data:         data,
			TemplateName: srcDirMeta.Name(),
		}

		srcDirPath := path.Join(config.Config.TemplatedSrcDir, config.TemplateName)
		srcDir := os.DirFS(srcDirPath)
		err := fs.WalkDir(srcDir, ".", config.buildWalk)
		check(err)
	}
}

func (config *buildConfig) buildWalk(path_ string, d fs.DirEntry, err error) error {
	excludeMatch, excludeIndex := PathMatch(path_, config.ExcludePaths)
	if excludeMatch != nil {
		logger.Debugf("Excluding path (%v): rule %v", *excludeMatch, excludeIndex+1)

		if d.IsDir() {
			return fs.SkipDir
		} else {
			return nil
		}
	}

	srcPath := path.Join(config.Config.TemplatedSrcDir, config.TemplateName, path_)
	dstPath := path.Join(config.DstDir, path_)

	if d.IsDir() {
		check(err)
		logger.Debugf("Building directory (%v -> %v)", srcPath, dstPath)

		err := os.Mkdir(dstPath, config.Config.DstMode.FileMode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			check(err)
		}
	} else {
		logger.Debugf("Building file (%v -> %v)", srcPath, dstPath)

		content, err := os.ReadFile(srcPath)
		check(err)

		contentTmpl, err := createTemplate("Content").Parse(string(content))
		if err != nil {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Failed to parse template (%v): %v", srcPath, err.Error())
			logger.SetOutput(os.Stdout)

			return nil
		}

		// Clone the template so that the original is not affected
		// This is used because if Content is entirely whitespace and comments,
		// the previous one is used, which we don't want
		newTmpl, err := config.Template.Clone()
		check(err)

		newTmpl.AddParseTree("Content", contentTmpl.Tree)

		dst, err := os.Create(dstPath)
		check(err)
		defer dst.Close()

		err = dst.Chmod(config.Config.DstMode.FileMode)
		check(err)

		data := BuildData(config, path_)

		if config.Config.FmtTemplatedHTML && path.Ext(d.Name()) == ".html" {
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

	return nil

}
