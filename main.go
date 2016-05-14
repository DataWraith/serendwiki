package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anknown/ahocorasick"
	"github.com/termie/go-shutil"
)

func init() {
	flag.Parse()
}

func printUsage() {
	fmt.Println("Usage: serendwiki <input-directory> <output-directory>")
}

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

func buildArticleMachine(fileList []string) goahocorasick.Machine {
	var (
		dict    [][]rune
		machine goahocorasick.Machine
	)

	for _, fn := range fileList {
		dict = append(dict, []rune(strings.ToLower(strings.TrimSpace(fn))))
	}

	if err := machine.Build(dict); err != nil {
		log.Fatalf("Error while building article recognizer: %s", err)
	}

	return machine
}

func processFiles(inputDir string, outputDir string, recognizer goahocorasick.Machine, linkTable map[string]string) int {
	var numArticles int

	copyOptions := &shutil.CopyTreeOptions{
		Symlinks:               true,
		Ignore:                 nil,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: true,
	}

	fileInfos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Error while reading input directory: %s", err)
	}

	for _, fi := range fileInfos {
		if isHiddenFile(fi.Name()) {
			continue
		}

		if fi.IsDir() {
			// Skip the output directory if it is inside the input directory
			if filepath.Join(inputDir, fi.Name()) == outputDir {
				continue
			}

			shutil.CopyTree(filepath.Join(inputDir, fi.Name()), filepath.Join(outputDir, fi.Name()), copyOptions)
			continue
		}

		if !isWikiFile(fi.Name()) {
			shutil.Copy(filepath.Join(inputDir, fi.Name()), outputDir, false)
			continue
		}

		processWikiFile(inputDir, outputDir, fi.Name(), fi.Mode(), recognizer, linkTable)
		numArticles++
	}

	return numArticles
}

func generateLinkTable(fileList []string) map[string]string {
	result := make(map[string]string)

	for _, fn := range fileList {
		result[strings.ToLower(fn)] = fn
	}

	return result
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
	err := os.Mkdir(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error: could not create output directory. Reason: %s", err)
	}

	fileList := gatherWikiFiles(inputDir)
	linkTable := generateLinkTable(fileList)
	recognizer := buildArticleMachine(fileList)

	numArticles := processFiles(inputDir, outputDir, recognizer, linkTable)

	fmt.Printf("Done processing wiki '%s', %d articles converted to HTML.\n", inputDir, numArticles)
}
