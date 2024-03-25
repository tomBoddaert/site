package main

import (
	"errors"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func createTemplate(name string) *template.Template {
	return template.New(name).Funcs(sprig.TxtFuncMap())
}

func GetTemplate(config *Config, name string) *template.Template {
	path_ := path.Join(config.TemplateDir, name)

	meta, err := os.Stat(path_)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Template does not exist (%v)", path_)
			logger.SetOutput(os.Stdout)

			return nil
		}

		check(err)
	}

	if meta.IsDir() {
		rootPath := path.Join(path_, config.TemplateRootFile)
		meta, err := os.Stat(rootPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				logger.SetOutput(os.Stderr)
				logger.Errorf("Template does not have a root file (%v)", rootPath)
				logger.SetOutput(os.Stdout)

				return nil
			}
			check(err)
		}
		if meta.IsDir() {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Template root is a directory (%v)", rootPath)
			logger.SetOutput(os.Stdout)

			return nil
		}

		logger.Debugf("Found template root (%v)", rootPath)

		paths := []string{}

		files, err := os.ReadDir(path_)
		check(err)
		for _, file := range files {
			filePath := path.Join(path_, file.Name())
			if file.IsDir() {
				logger.SetOutput(os.Stderr)
				logger.Errorf("Template contains a directory (%v)", filePath)
				logger.SetOutput(os.Stdout)

				return nil
			}

			paths = append(paths, filePath)

			logger.Debugf("Found template file (%v)", filePath)
		}

		tmpl, err := createTemplate(config.TemplateRootFile).
			ParseFiles(paths...)
		if err != nil {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Template contains errors (%v): %v", name, err.Error())
			logger.SetOutput(os.Stdout)

			return nil
		}

		return tmpl
	} else {
		tmpl, err := createTemplate(name).
			ParseFiles(path_)

		logger.Debugf("Found single-file template (%v)", path_)

		if err != nil {
			logger.SetOutput(os.Stderr)
			logger.Errorf("Template contains errors (%v): %v", name, err.Error())
			logger.SetOutput(os.Stdout)

			return nil
		}

		return tmpl
	}
}
