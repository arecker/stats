package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"sync"
)

func processEntry(entry *Entry, wg *sync.WaitGroup) {
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
		log.Fatal(error)
	}

	// write to file
	if error = ioutil.WriteFile(entry.WordsFilePath, data, 0644); error != nil {
		log.Fatal(error)
	}
}

func main() {
	var wg sync.WaitGroup

	for _, entry := range Entries() {
		wg.Add(1)
		go processEntry(entry, &wg)
	}

	wg.Wait()
}
