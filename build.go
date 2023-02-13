package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/yosssi/gohtml"
)

var includeTSFiles = false
var formatHTMLFiles = false

var tempDir string

func build() {
	// Create the temp directory
	var err error
	// The os temp directory is not used because some OSs
	//  don't support moving files between volumes / devices
	// tempDir, err = os.MkdirTemp(os.TempDir(), "site-docs")
	tempDir, err = os.MkdirTemp(".", "site-docs-new-")
	check(err)

	// Settings
	gohtml.Condense = true

	copyRawPages("")
	copyTemplatedPages()

	// Move temp directory to docs directory
	check(os.RemoveAll("docs"))
	check(os.Rename(tempDir, "docs"))

	transpileTS()
}

func copyRawPages(subDir string) {
	// Read raw pages
	rawPages, err := os.ReadDir(path.Join("rawPages", subDir))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println("Warning: rawPages directory does not exist")
		return
	}
	check(err)

	// Loop over raw pages
	for _, rawPage := range rawPages {
		// If it is a directory, run this function over that too
		if rawPage.Type().IsDir() {
			err = os.Mkdir(path.Join(tempDir, subDir, rawPage.Name()), 0755)
			if err != nil && !errors.Is(err, os.ErrExist) {
				check(err)
			}

			copyRawPages(path.Join(subDir, rawPage.Name()))
			continue

		} else if !rawPage.Type().IsRegular() {
			fmt.Printf("Not a regular file or directory: %s", path.Join("rawPages", subDir, rawPage.Name()))
			continue
		}

		// If it is a TypeScript file and 'inclts' was not provided, skip it
		if !includeTSFiles && strings.HasSuffix(rawPage.Name(), ".ts") {
			continue
		}

		// Copy the raw file to docs
		rawPageFile, err := os.ReadFile(path.Join("rawPages", subDir, rawPage.Name()))
		check(err)

		// If the file is a html file and 'fmthtml' was provided, format it
		if formatHTMLFiles && strings.HasSuffix(rawPage.Name(), ".html") {
			rawPageFile = gohtml.FormatBytes(rawPageFile)
		}

		check(os.WriteFile(
			path.Join(tempDir, subDir, rawPage.Name()),
			rawPageFile,
			0644,
		))
	}
}

func copyTemplatedPages() {
	// Get page variables
	pageVariables := map[string]map[string]any{}
	pageVariablesFile, err := os.ReadFile("pageVariables.json")
	if !(err == nil || errors.Is(err, os.ErrNotExist)) {
		check(err)
	}
	if err == nil {
		json.Unmarshal(pageVariablesFile, &pageVariables)
	}

	// Get templates
	templateFiles, err := os.ReadDir("templates")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println("Warning: templates directory does not exist")
		templateFiles = make([]fs.DirEntry, 0)
	} else {
		check(err)
	}

	// Loop over templates
	for _, templateFile := range templateFiles {
		// Template must be built with the Parse(string) method to add all the files to the same template
		//  Using ParseFiles or ParseGlob creates a new template for each file
		tmpl := template.New("template")

		// If the template is split over files, add all the files
		if templateFile.Type().IsDir() {
			templateDir, err := os.ReadDir(path.Join("templates", templateFile.Name()))
			check(err)

			for _, templatePartFile := range templateDir {
				if !templatePartFile.Type().IsRegular() {
					fmt.Printf("Not a regular file: %s", path.Join("templates", templateFile.Name(), templatePartFile.Name()))
					continue
				}

				// Read the file
				templatePart, err := os.ReadFile(path.Join("templates", templateFile.Name(), templatePartFile.Name()))
				check(err)

				// Add it to the template
				tmpl, err = tmpl.Parse(string(templatePart))
				check(err)
			}

		} else if templateFile.Type().IsRegular() {
			// Read the file
			templatePart, err := os.ReadFile(path.Join("templates", templateFile.Name()))
			check(err)

			// Add it to the template
			tmpl, err = tmpl.Parse(string(templatePart))
			check(err)

		} else {
			fmt.Printf("Not a regular file or directory: %s", path.Join("templates", templateFile.Name()))
			continue
		}

		copyTemplatedDir(templateFile.Name(), "", tmpl, pageVariables)
	}
}

func copyTemplatedDir(templateName string, subDir string, tmpl *template.Template, pageVariables map[string]map[string]any) {
	// Read templated pages
	templatedPages, err := os.ReadDir(path.Join("templatedPages", templateName, subDir))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Warning: no templatedPages directory for template \"%s\"\n", path.Join(templateName, subDir))
		return
	}
	check(err)

	// Loop over templated pages
	for _, templatedPage := range templatedPages {
		// If it is a directory, run this function over that too
		if templatedPage.Type().IsDir() {
			err = os.Mkdir(path.Join(tempDir, subDir, templatedPage.Name()), 0755)
			if err != nil && !errors.Is(err, os.ErrExist) {
				check(err)
			}

			copyTemplatedDir(templateName, path.Join(subDir, templatedPage.Name()), tmpl, pageVariables)
			continue

		} else if !templatedPage.Type().IsRegular() {
			fmt.Printf("Not a regular file or directory: %s", path.Join("templatedPages", templateName, subDir, templatedPage.Name()))
			continue
		}

		// Read content of file
		contentFile, err := os.ReadFile(path.Join("templatedPages", templateName, subDir, templatedPage.Name()))
		check(err)

		pageTmpl, err := tmpl.Clone()
		check(err)

		// Add the file to the template as 'Content'
		pageTmpl, err = pageTmpl.New("Content").Parse(string(contentFile))
		check(err)

		// Get the page variables
		data := map[string]any{}
		if len(pageVariables["default"]) != 0 {
			for key, value := range pageVariables["default"] {
				data[key] = value
			}
		}
		if len(pageVariables[templateName]) != 0 {
			for key, value := range pageVariables[templateName] {
				data[key] = value
			}
		}
		templatePath := path.Join(templateName, subDir, templatedPage.Name())
		if len(pageVariables[templatePath]) != 0 {
			for key, value := range pageVariables[templatePath] {
				data[key] = value
			}
		}

		pageBuffer := bytes.NewBuffer(nil)
		check(pageTmpl.ExecuteTemplate(pageBuffer, "template", data))

		// Reformat
		//  Not using gohtml writer because leading blank lines result in a blank output
		formattedBuffer := gohtml.FormatBytes(pageBuffer.Bytes())

		// Write page to docs
		check(os.WriteFile(
			path.Join(tempDir, subDir, templatedPage.Name()),
			formattedBuffer,
			0644,
		))
	}
}

func transpileTS() {
	// Transpile the TypeScript files
	cmd := exec.Command("npx", "tsc")
	output := bytes.NewBuffer(nil)
	cmd.Stdout = output

	err := cmd.Run()
	if err != nil {
		// If the error is TS18003, no files were found to transpile, so skip
		noFiles := strings.Contains(output.String(), "error TS18003")
		if noFiles {
			return
		}

		noTsc := strings.Contains(output.String(), "This is not the tsc command you are looking for")
		if noTsc {
			fmt.Println("TypeScript not installed!")
			fmt.Println("Install typescript using 'npm i typescript'")
			fmt.Println("TypeScript files not transpiled!")
			return
		}

		fmt.Println(output.String())
		fmt.Println("TypeScript files not transpiled!")
		return
	}

	// If there is no error and an output, print it
	if output.Len() != 0 {
		fmt.Println("From typescript: (npx tsc)")
		fmt.Println(output.String())
	}
}
