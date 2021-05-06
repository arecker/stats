package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
)

const EntriesDir string = `blog/_posts`
const WordsDir string = `words/`

type Entry struct {
	FileName      string
	FilePath      string
	DateSlug      string
	WordsFilePath string
}

func NewEntry(fileName string) *Entry {
	entry := new(Entry)
	entry.FileName = fileName
	entry.DateSlug = fileName[0:10]
	entry.FilePath = path.Join(EntriesDir, fileName)
	entry.WordsFilePath = path.Join(WordsDir, entry.DateSlug+".json")
	return entry
}

func Entries() []*Entry {
	var entries []*Entry

	results, error := ioutil.ReadDir(EntriesDir)
	if error != nil {
		log.Fatal(error)
	}

	for _, result := range results {
		if !result.IsDir() {
			entries = append(entries, NewEntry(result.Name()))
		}
	}

	return entries
}

func (entry *Entry) rawContent() (string, error) {
	result, error := ioutil.ReadFile(entry.FilePath)

	if error != nil {
		return "", error
	}

	return string(result), error
}

func (entry *Entry) content() (string, error) {
	frontmatterPattern := regexp.MustCompile(`(?s)^---\n.*?\n---\n`)
	raw, error := entry.rawContent()
	if error != nil {
		return raw, error
	}

	return frontmatterPattern.ReplaceAllString(raw, ""), nil
}

func (entry *Entry) wordScanner() (*bufio.Scanner, error) {
	var scanner *bufio.Scanner

	content, error := entry.content()
	if error != nil {
		return scanner, error
	}

	reader := strings.NewReader(content)
	scanner = bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	return scanner, error
}

func sanitizeWord(word string) string {
	// reject opening and closing HTML tags
	if regexp.MustCompile(`^<.*|\/>$`).MatchString(word) {
		return ""
	}

	// reject URLs
	if regexp.MustCompile(`^https?:\/\/`).MatchString(word) {
		return ""
	}

	// reject HTML attributes
	if regexp.MustCompile(`^.*=`).MatchString(word) {
		return ""
	}

	// reject markdown horizontal rules
	if regexp.MustCompile(`^[\-]{2,3}$`).MatchString(word) {
		return ""
	}

	// reject paths
	if regexp.MustCompile(`^(\/.*?)+$`).MatchString(word) {
		return ""
	}

	word = strings.ToLower(word)

	// strip out illegal characters
	word = regexp.MustCompile(`[~“”‘’–…\/\{\}\[\]\\\*#()"_?!,.]`).ReplaceAllString(word, "")

	// strip out back ticks (nice one, go)
	word = regexp.MustCompile("`").ReplaceAllString(word, "")

	// reject times
	if regexp.MustCompile(`^[0-9]{1,2}[:]?[0-9]{0,2}(am|pm|AM|PM)?$`).MatchString(word) {
		return ""
	}

	// wrapped single quotes, e.g. "He said 'hello',"
	word = regexp.MustCompile(`^'(.*?)'$`).ReplaceAllString(word, "$1")

	// leading single quotes, e.g. 'Hello
	word = regexp.MustCompile(`^'(.*?)$`).ReplaceAllString(word, "$1")

	// check for dangling symbols
	word = regexp.MustCompile(`^[\|à§\/\-\+%–&<>]$`).ReplaceAllString(word, "")

	// check for dangling numbers
	word = regexp.MustCompile(`^\d+$`).ReplaceAllString(word, "")

	return word
}

func (entry *Entry) Words() ([]string, error) {
	var words []string

	scanner, error := entry.wordScanner()
	if error != nil {
		return words, error
	}

	for scanner.Scan() {
		result := scanner.Text()
		result = sanitizeWord(result)
		if result != "" {
			words = append(words, result)
		}
	}

	return words, error
}
