package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/anknown/ahocorasick"
	"github.com/russross/blackfriday"
)

func processWikiFile(inputDir string, outputDir string, fileName string, mode os.FileMode, recognizer goahocorasick.Machine) {
	contents, err := ioutil.ReadFile(filepath.Join(inputDir, fileName))
	if err != nil {
		log.Fatalf("Error while reading '%s': %s", fileName, err)
	}

	output := blackfriday.MarkdownCommon(contents)

	// TODO: Link articles

	outPath := filepath.Join(outputDir, fileName+".html")
	err = ioutil.WriteFile(outPath, []byte(output), mode)
	if err != nil {
		log.Fatalf("Error while writing to '%s': %s", outPath, err)
	}
}
