package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func init() {
	flag.Parse()
}

func printUsage() {
	fmt.Println("Usage: serendwiki <input-directory> <output-directory>")
}

func gatherFilesWithExt(inputDir string, ext string) []string {
	matches, err := filepath.Glob(inputDir + "/*." + ext)
	if err != nil {
		log.Fatalf("Error while searching for markdown files: %s", err)
	}
	return matches
}

func gatherFiles(inputDir string) []string {
	var matches []string

	matches = gatherFilesWithExt(inputDir, "[Mm][Dd]")
	matches = append(matches, gatherFilesWithExt(inputDir, "[Mm][Kk][Dd]")...)
	matches = append(matches, gatherFilesWithExt(inputDir, "[Mm][Kk][Dd][Nn]")...)
	matches = append(matches, gatherFilesWithExt(inputDir, "[Mm][Aa][Rr][Kk][Dd][Oo][Ww][Nn]")...)

	if len(matches) == 0 {
		log.Fatalf("No markdown files (*.md, *.mkd, *.mkdn, *.markdown) found in %s", inputDir)
	}

	return matches
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
	if len(flag.Args()) < 2 {
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

	fileList := gatherFiles(inputDir)
	fmt.Println(fileList)
}
