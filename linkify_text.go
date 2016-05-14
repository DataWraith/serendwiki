package main

import (
	"bytes"
	"sort"
	"strings"
	"unicode"

	"github.com/anknown/ahocorasick"
)

// term holds an article title, the position of that title in the input and a
// priority for that article title. Articles with lower priority take precedence
// over higher priority titles when article titles overlap.
//
type term struct {
	Title     []rune
	Pos      int
	Priority int
}

type byLength []term // byLength is a helper-type for sorting terms by title length
type byPos []term // byPos is a helper-type for sorting terms by title position

func (t byPos) Len() int           { return len(t) }
func (t byPos) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byPos) Less(i, j int) bool { return t[i].Pos < t[j].Pos }

func (t byLength) Len() int           { return len(t) }
func (t byLength) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byLength) Less(i, j int) bool { return len(t[i].Title) > len(t[j].Title) }

// removeOverlap takes a list of titles (and their positions) that were matched
// in the input and removes any overlap of titles by throwing out shorter titles
// and keeping longer, more specific titles if they happen to overlap.
//
func removeOverlap(input []*goahocorasick.Term) []term {
	terms := make([]term, 0, len(input))

	for _, t := range input {
		terms = append(terms, term{Title: t.Word, Pos: t.Pos})
	}

	if len(input) < 2 {
		return terms
	}

	sort.Sort(byLength(terms))

	for i := range terms {
		terms[i].Priority = i
	}

	sort.Sort(byPos(terms))

	result := make([]term, 0, len(input))

	i := 0
	j := 1

	for i < len(input) {
		if j >= len(input) {
			result = append(result, terms[i])
			break
		}

		if terms[i].Pos+len(terms[i].Title) <= terms[j].Pos {
			result = append(result, terms[i])
			i = j
			j++
			continue
		}

		if terms[i].Priority < terms[j].Priority {
			j++
		} else {
			i = j
			j++
		}
	}

	return result
}

// linkifyText uses the Aho Corasick recognizer and the provided linkTable to
// replace HTML text with links to the recognized article titles.
//
func linkifyText(input []byte, recognizer goahocorasick.Machine, linkTable map[string]string) []byte {
	rInput := bytes.Runes(input)
	rInputLower := make([]rune, 0, len(rInput))
	for i := 0; i < len(rInput); i++ {
		rInputLower = append(rInputLower, unicode.ToLower(rInput[i]))
	}

	searchResults := recognizer.MultiPatternSearch(rInputLower, false)
	terms := removeOverlap(searchResults)
	terms = append(terms, term{Pos: len(rInput), Title: []rune{}, Priority: 1 << 30})

	curTerm := 0
	rresult := []rune{}

	i := 0
	for i < len(rInput) {
		if i < terms[curTerm].Pos {
			rresult = append(rresult, rInput[i])
			i++
			continue
		}

		rresult = append(rresult, []rune("<a href=\"")...)
		rresult = append(rresult, []rune(linkTable[strings.ToLower(string(rInput[i:i+len(terms[curTerm].Title)]))])...)
		rresult = append(rresult, []rune(".html\">")...)
		rresult = append(rresult, rInput[i:i+len(terms[curTerm].Title)]...)
		rresult = append(rresult, []rune("</a>")...)

		i += len(terms[curTerm].Title)
		curTerm++
	}

	return []byte(string(rresult))
}
