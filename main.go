package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error! Expected exactly 2 argument: url and output path")
		return
	}

	website := os.Args[1]
	outputPath := os.Args[2]

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err = os.Mkdir(outputPath, os.ModeDir)
		if err != nil {
			fmt.Println("Failed to create output path", err)
		} else {
			fmt.Println("Created output path:", outputPath)
		}
	}

	websiteBytes, err := downloadBytes(website)

	if err != nil {
		fmt.Println("Failed to download website", err)
	}

	websiteString := string(websiteBytes)
	webSources, err := findAllWebSources(websiteString)

	if err != nil {
		fmt.Println("Failed to find sources", err)
	}

	fmt.Println("Found", len(webSources), "sources")

	for _, source := range webSources {
		bytes, _ := downloadBytes(source)

		if len(bytes) == 0 {
			fmt.Println("Empty resource:", source)
			continue
		}

		nameStart := strings.LastIndex(source, "/")
		newName := source[nameStart+1:]
		queryStart := strings.Index(newName, "?")

		if queryStart >= 0 {
			newName = newName[:queryStart]
		}

		filePAth := filepath.Join(outputPath, newName)
		err = writeBytesToFile(bytes, filePAth)

		if err != nil {
			fmt.Println("Failed to write file", filePAth, err)
		} else {
			fmt.Println("Wrote file", filePAth)
		}
	}
}

func findAllWebSources(html string) ([]string, error) {
	re, err := regexp.Compile("src=\"(.*?)\"")

	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(html, -1)

	sources := make([]string, len(matches))

	for i, match := range matches {
		sources[i] = match[1]
	}

	return sources, nil
}

func downloadBytes(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

func writeBytesToFile(image []byte, filename string) error {
	return ioutil.WriteFile(filename, image, 0644)
}
