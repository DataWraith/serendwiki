package main

import (
	"io/ioutil"
	"log"
	"strings"
)

// isWikiFile returns whether a filename is considered a wiki article
func isWikiFile(fileName string) bool {
	return !strings.Contains(fileName, ".")
}

// isHiddenFile returns whether a file is a dotfile (hidden)
//
// There is no Windows support currently.
//
func isHiddenFile(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}

// gatherWikiFiles scans the input directory and returns a list of filenames
// that appear to be wiki articles (don't have file extensions)
//
func gatherWikiFiles(inputDir string) []string {
	var wikiFiles []string

	fileInfos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Error while reading input directory: %s", err)
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}

		if isWikiFile(fi.Name()) {
			wikiFiles = append(wikiFiles, fi.Name())
		}
	}

	if len(wikiFiles) == 0 {
		log.Fatalf("No wiki files (files without extension) found in %s", inputDir)
	}

	return wikiFiles
}
