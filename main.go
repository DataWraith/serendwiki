package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	flag.Parse()
}

func printUsage() {
	fmt.Println("Usage: serendwiki <input-directory> <output-directory>")
}

func isWikiFile(fileName string) bool {
	if strings.Contains(fileName, ".") {
		return false
	}

	return true
}

func isHiddenFile(fileName string) bool {
	if strings.HasPrefix(fileName, ".") {
		return true
	}

	return false
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

func checkForErrors(inputDir string, outputDir string) {
	if _, err := os.Stat(inputDir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Input directory '%s' does not exist.", inputDir)
		}
		log.Fatalf("Error while opening input directory '%s': %s", inputDir, err)
	}

	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		log.Fatalf("Refusing to overwrite exisitng directory '%s'", outputDir)
	}
}

func main() {
	if len(flag.Args()) != 2 {
		printUsage()
		os.Exit(1)
	}

	inputDir := filepath.Clean(flag.Args()[0])
	outputDir := filepath.Clean(flag.Args()[1])

	checkForErrors(inputDir, outputDir)

	// Create output directory
	err := os.Mkdir(outputDir, 0644)
	if err != nil {
		log.Fatalf("Error: could not create output directory. Reason: %s", err)
	}

	fileList := gatherWikiFiles(inputDir)
	fmt.Println(fileList)

	// TODO: Copy all files that are not wiki files from the input directory to the output directory (except '.git', etc.)
}
