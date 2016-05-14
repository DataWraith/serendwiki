package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

type term struct {
	Word     []rune
	Pos      int
	Priority int
}

type ByLength []term
type ByPos []term

func (t ByPos) Len() int           { return len(t) }
func (t ByPos) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByPos) Less(i, j int) bool { return t[i].Pos < t[j].Pos }

func (t ByLength) Len() int           { return len(t) }
func (t ByLength) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByLength) Less(i, j int) bool { return len(t[i].Word) > len(t[j].Word) }

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("article").Parse(defaultTemplate))
}

func removeOverlap(input []*goahocorasick.Term) []term {
	terms := make([]term, 0, len(input))

	for _, t := range input {
		terms = append(terms, term{Word: t.Word, Pos: t.Pos})
	}

	sort.Sort(ByLength(terms))

	for i := range terms {
		terms[i].Priority = i
	}

	sort.Sort(ByPos(terms))

	result := make([]term, 0, len(input))

	i := 0
	j := 1

	if len(input) < 2 {
		return terms
	}

	appendLast := false
	for j < len(input) {
		if terms[i].Pos+len(terms[i].Word) < terms[j].Pos {
			result = append(result, terms[i])
			appendLast = false
			i = j
			j++
			continue
		}

		appendLast = true

		if terms[i].Priority < terms[j].Priority {
			j++
		} else {
			i = j
			j++
		}
	}

	if appendLast {
		result = append(result, terms[i])
	}

	return result
}

func linkifyText(input []byte, recognizer goahocorasick.Machine, linkTable map[string]string) []byte {
	rinput := bytes.Runes(input)
	rinput_lower := make([]rune, 0, len(rinput))
	for i := 0; i < len(rinput); i++ {
		rinput_lower = append(rinput_lower, unicode.ToLower(rinput[i]))
	}

	searchResults := recognizer.MultiPatternSearch(rinput_lower, false)
	terms := removeOverlap(searchResults)
	terms = append(terms, term{Pos: len(rinput), Word: []rune{}, Priority: ^0})

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
		rresult = append(rresult, []rune(linkTable[strings.ToLower(string(rinput[i:i+len(terms[curTerm].Word)]))])...)
		rresult = append(rresult, []rune(".html\">")...)
		rresult = append(rresult, rinput[i:i+len(terms[curTerm].Word)]...)
		rresult = append(rresult, []rune("</a>")...)

		i += len(terms[curTerm].Word)
		curTerm++
	}

	return []byte(string(rresult))
}

func linkArticles(inputHtml []byte, recognizer goahocorasick.Machine, linkTable map[string]string) []byte {
	z := html.NewTokenizer(bytes.NewReader(inputHtml))
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
