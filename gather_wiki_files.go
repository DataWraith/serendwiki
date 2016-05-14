package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func isWikiFile(fileName string) bool {
	return !strings.Contains(fileName, ".")
}

func isHiddenFile(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}

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
