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

// linkArticles takes the HTML output of Blackfriday and parses it to figure out
// where text that needs to be "linkified" is.
//
func linkArticles(inputHTML []byte, recognizer goahocorasick.Machine, linkTable map[string]string) []byte {
	z := html.NewTokenizer(bytes.NewReader(inputHTML))
	result := []byte{}

	insideLinkTag := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return result
			}
			log.Fatalf("Error while parsing generated HTML: %s", z.Err())

		case html.TextToken:
			if !insideLinkTag {
				result = append(result, linkifyText(z.Text(), recognizer, linkTable)...)
			} else {
				result = append(result, z.Raw()...)
			}

		case html.StartTagToken:
			tn, _ := z.TagName()
			if string(tn) == "a" || string(tn) == "A" {
				insideLinkTag = true
			}
			result = append(result, z.Raw()...)

		case html.EndTagToken:
			tn, _ := z.TagName()
			if string(tn) == "a" || string(tn) == "A" {
				insideLinkTag = false
			}
			result = append(result, z.Raw()...)

		default:
			result = append(result, z.Raw()...)
		}
	}
}

// processWikiFile takes the filename of an article (amongs other things),
// converts the markdown markup into HTML, "linkifies" said HTML and then sanitizes
// it using Bluemonday to prevent malicious Javascript and the like.
//
func processWikiFile(inputDir string, outputDir string, fileName string, mode os.FileMode, recognizer goahocorasick.Machine, linkTable map[string]string) {
	contents, err := ioutil.ReadFile(filepath.Join(inputDir, fileName))
	if err != nil {
		log.Fatalf("Error while reading '%s': %s", fileName, err)
	}

	output := blackfriday.MarkdownCommon(contents)
	linked := linkArticles(output, recognizer, linkTable)
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
