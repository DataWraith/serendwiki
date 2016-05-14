package main

import (
	"text/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/anknown/ahocorasick"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const (
	defaultTemplate = `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
			</head>
			<body>
				{{.}}
			</body>
		</html>
	`
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("article").Parse(defaultTemplate))
}

func linkArticles(inputHtml []byte, recognizer goahocorasick.Machine) []byte {
	return inputHtml
}

func processWikiFile(inputDir string, outputDir string, fileName string, mode os.FileMode, recognizer goahocorasick.Machine) {
	contents, err := ioutil.ReadFile(filepath.Join(inputDir, fileName))
	if err != nil {
		log.Fatalf("Error while reading '%s': %s", fileName, err)
	}

	output := blackfriday.MarkdownCommon(contents)
	linked := linkArticles(output, recognizer)
	safe := bluemonday.UGCPolicy().SanitizeBytes(linked)

	outPath := filepath.Join(outputDir, fileName+".html")

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Error creating file '%s': %s", outPath, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, string(safe))
	if err != nil {
		log.Fatalf("Error while writing to '%s': %s", outPath, err)
	}
}
