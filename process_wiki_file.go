package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/anknown/ahocorasick"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"golang.org/x/net/html"
)

const (
	defaultTemplate = `<!DOCTYPE html><html><head><meta charset="utf-8"></head><body>{{.}}</body></html> `
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("article").Parse(defaultTemplate))
}

func linkifyText(input []byte) []byte {
	return input
}

func linkArticles(inputHtml []byte, recognizer goahocorasick.Machine) []byte {
	z := html.NewTokenizer(bytes.NewReader(inputHtml))
	result := []byte{}

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return result
			}
			log.Fatalf("Error while parsing generated HTML: %s", z.Err())

		case html.TextToken:
			result = append(result, linkifyText(z.Text())...)

		default:
			result = append(result, z.Raw()...)
		}
	}
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
