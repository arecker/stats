package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"sync"

	progressbar "github.com/schollz/progressbar/v3"
)

func processEntry(entry *Entry, wg *sync.WaitGroup, progressBar *progressbar.ProgressBar) {
	defer wg.Done()

	// read in words
	words, error := entry.Words()
	if error != nil {
		log.Fatal(error)
	}

	// sort and encode
	sort.Strings(words)
	data, err := json.MarshalIndent(words, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// write to file
	if error = ioutil.WriteFile(entry.WordsFilePath, data, 0644); error != nil {
		log.Fatal(error)
	}

	progressBar.Add(1)
}

func main() {
	var wg sync.WaitGroup
	entries := Entries()
	progressBar := progressbar.Default(int64(len(entries)), "Scanning entries")

	for _, entry := range Entries() {
		wg.Add(1)
		go processEntry(entry, &wg, progressBar)
	}

	wg.Wait()

	progressBar = progressbar.Default(int64(len(entries)), "Compiling total words")
	var allWords []string
	for _, entry := range entries {
		reader, error := os.Open(entry.WordsFilePath)
		if error != nil {
			log.Fatal(error)
		}

		byteValue, error := ioutil.ReadAll(reader)
		if error != nil {
			log.Fatal(error)
		}

		var theseWords []string
		json.Unmarshal(byteValue, &theseWords)
		allWords = append(allWords, theseWords...)
		progressBar.Add(1)
	}

	fmt.Printf("Building all words (%d)\n", len(allWords))
	sort.Strings(allWords)
	data, error := json.MarshalIndent(allWords, "", "  ")
	if error != nil {
		log.Fatal(error)
	}

	if error = ioutil.WriteFile(path.Join(WordsDir, "all.json"), data, 0644); error != nil {
		log.Fatal(error)
	}

	fmt.Println("Done!")
}
