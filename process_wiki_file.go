package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/anknown/ahocorasick"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

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
	err = ioutil.WriteFile(outPath, safe, mode)
	if err != nil {
		log.Fatalf("Error while writing to '%s': %s", outPath, err)
	}
}
