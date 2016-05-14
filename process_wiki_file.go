package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"unicode"

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

func removeOverlap(input []*goahocorasick.Term) []*goahocorasick.Term {
	return input
}

func linkifyText(input []byte, recognizer goahocorasick.Machine) []byte {
	rinput := bytes.Runes(input)
	rinput_lower := make([]rune, 0, len(rinput))
	for i := 0; i < len(rinput); i++ {
		rinput_lower = append(rinput_lower, unicode.ToLower(rinput[i]))
	}

	terms := recognizer.MultiPatternSearch(rinput_lower, false)
	terms = removeOverlap(terms)
	terms = append(terms, &goahocorasick.Term{Pos: len(rinput), Word: []rune{}})

	curTerm := 0
	rresult := []rune{}

	i := 0
	for i < len(rinput) {
		if i < terms[curTerm].Pos {
			rresult = append(rresult, rinput[i])
			i++
			continue
		}

		rresult = append(rresult, []rune("<a href=\"")...)
		rresult = append(rresult, rinput[i:i+len(terms[curTerm].Word)]...)
		rresult = append(rresult, []rune(".html\">")...)
		rresult = append(rresult, rinput[i:i+len(terms[curTerm].Word)]...)
		rresult = append(rresult, []rune("</a>")...)

		i += len(terms[curTerm].Word)
		curTerm++
	}

	return []byte(string(rresult))
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
			result = append(result, linkifyText(z.Text(), recognizer)...)

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
