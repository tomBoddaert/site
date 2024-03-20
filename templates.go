package main

import (
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func appendFileToBuf(name string, dst io.Writer) {
	file, err := os.Open(name)
	check(err)
	defer file.Close()

	_, err = io.Copy(dst, file)
	check(err)
}

func getTemplateText(name string) string {
	meta, err := os.Stat(name)
	check(err)

	buf := new(strings.Builder)

	if meta.IsDir() {
		dir, err := os.ReadDir(name)
		check(err)

		for _, fileMeta := range dir {
			if fileMeta.Type().IsDir() {
				logger.Warnf("%v is a directory! Ignoring", path.Join(name, fileMeta.Name()))
				continue
			}

			filePath := path.Join(name, fileMeta.Name())
			appendFileToBuf(filePath, buf)
		}
	} else {
		appendFileToBuf(name, buf)
	}

	return buf.String()
}

func createTemplate(name string) *template.Template {
	return template.New(name).Funcs(sprig.TxtFuncMap())
}

func getTemplate(name string) *template.Template {
	tmplText := getTemplateText(name)
	tmpl, err := createTemplate("template").Parse(tmplText)
	check(err)

	return tmpl
}
