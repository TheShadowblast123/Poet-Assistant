package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type Stanza struct {
	lines    []string
	endWords []string
}
type JsonOrError struct {
	IsJson bool
	result []struct {
		Word string `json:"word"`
	}
	e error
}

const baseDatamuseUrl string = `https://api.datamuse.com/words`

type RhymesData struct {
	Rhymes   [][][]string `json:"rhymes"`
	Synonyms [][]string   `json:"synonyms"`
}

var rhymeGroups [][][]string
var synonymGroups [][]string

func loadJSONData() error {
	file, err := os.Open("data.json")
	if err != nil {
		return err

	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var rhymesData RhymesData
	if err := json.Unmarshal(data, &rhymesData); err != nil {
		return err
	}

	rhymeGroups = rhymesData.Rhymes
	synonymGroups = rhymesData.Synonyms
	return nil
}

func findRhymeGroup(word string) []string {
	order := []int{0, 1, 2, 3, 7, 4, 6, 5, 9}
	var index int = -1
	var answer []string
	var length int

	if len(word) < 9 {
		length = len(word)
	} else {
		length = 0
	}

	for _, group := range rhymeGroups[length] {
		x := indexOf(group, word)
		if x != -1 {
			if x == 0 {
				answer = group
				index = 0
				break
			}
			if index == -1 || index > x {
				answer = group
				index = x
			}
		}
	}

	if index == 0 || index == 1 {
		return answer
	}

	for _, num := range order {
		if num == length {
			continue
		}
		for _, group := range rhymeGroups[num] {
			x := indexOf(group, word)
			if x != -1 {
				if x == 1 {
					answer = group
					index = 0
					break
				}
				if index == -1 || index > x {
					answer = group
					index = x
				}
			}
		}
	}

	if index == -1 {
		return []string{"Didn't find anything"}
	}

	return answer
}

func findRhymeAllGroups(word string) [][]string {
	order := []int{0, 1, 2, 3, 7, 4, 6, 5, 8}
	var answer [][]string

	for _, num := range order {
		for _, group := range rhymeGroups[num] {
			x := indexOf(group, word)
			if x != -1 {
				answer = append(answer, group)
			}
		}
	}

	return answer
}

func deepDoesRhyme(s, w string) bool {
	sRhymes := findRhymeAllGroups(s)
	wRhymes := findRhymeAllGroups(w)
	if len(sRhymes) > 0 {
		for _, group := range sRhymes {
			for _, word := range group {
				if doesRhyme(w, word) {
					return true
				}
			}
		}
	}
	if len(wRhymes) > 0 {
		for _, group := range wRhymes {
			for _, word := range group {
				if doesRhyme(s, word) {
					return true
				}
			}
		}
	}
	return false
}
func doesRhyme(s, w string) bool {
	if s == w {
		return true
	}
	sRhymes := findRhymeAllGroups(s)

	if len(sRhymes) == 1 && sRhymes[0][0] == "Didn't find anything" {
		wRhymes := findRhymeAllGroups(w)
		for _, group := range wRhymes {
			for _, rhymingWord := range group {
				if rhymingWord == s {
					return true
				}
			}
		}
		return false
	}

	for _, group := range sRhymes {
		for _, rhymingWord := range group {
			if rhymingWord == w {
				return true
			}
		}
	}
	wRhymes := findRhymeAllGroups(w)

	if len(wRhymes) == 1 && wRhymes[0][0] == "Didn't find anything" {
		return false
	}
	for _, group := range wRhymes {
		for _, rhymingWord := range group {
			if rhymingWord == s {
				return true
			}
		}
	}
	return false
}

func indexOf(slice []string, target string) int {
	for i, element := range slice {
		if element == target {
			return i
		}
	}
	return -1
}

func CountWordSyllables(w string) int {
	w = strings.ToLower(w)
	fmt.Println(w)
	// Handle exceptions
	twoSyllable := []string{"coapt", "coed", "coinci", "colonel", "cafe", "scotia"}
	if len(w) <= 3 || w == "preach" || w == "preyed" {
		return 1
	}
	if contains(twoSyllable, w) {
		return 2
	}
	if w == "serious" || w == "worcestershire" || w == "alias" || w == "acacia" {
		return 3
	}
	if w == "epitome" || w == "hyperbole" {
		return 4
	}

	regex := []*regexp.Regexp{
		regexp.MustCompile(`[aeiouy]{1,}`),
		regexp.MustCompile(`^(?:mc)`),
		regexp.MustCompile(`^(?:tri)[aeiouy]`),
		regexp.MustCompile(`^(?:bi)[aeiouy]`),
		regexp.MustCompile(`^(?:pre)[aeiouy]`),
		regexp.MustCompile(`[^tc](?:ian)$`),
		regexp.MustCompile(`[aeiou][^aeiouy]e$`),
		regexp.MustCompile(`[aeiou][^aeiouy]es$`),
		regexp.MustCompile(`i[ao]$`),
		regexp.MustCompile(`[aeiouy]ing$`),
	}

	// Base case
	holdVowels := regex[0].FindAllString(w, -1) // count vowel groups
	fmt.Println(holdVowels)
	count := len(holdVowels)
	if count == 0 {
		return 0
	}
	if regex[1].MatchString(w) || regex[2].MatchString(w) || regex[3].MatchString(w) || regex[4].MatchString(w) {
		count++ // Handle syllabic prefixes
	}
	if regex[5].MatchString(w) {
		count++ // Handle ian suffix
	}
	if regex[6].MatchString(w) || regex[7].MatchString(w) {
		count-- // Handle [vowel]_e(s)
	}
	if regex[8].MatchString(w) {
		count++ // Handle i[ao]
	}
	if regex[9].MatchString(w) {
		count++ // Handle vowel+ing
	}
	return count
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func localSyllables(str string) int {
	hold := strings.Split(str, " ")
	sum := 0
	if hold != nil {
		for _, word := range hold {
			sum += CountWordSyllables(removePunctuationFromWord(word))
		}
	}
	return sum
}
func syllableCount(str string) int {
	apiURL := fmt.Sprintf("%s?sp=%s&md=s", baseDatamuseUrl, str)

	res, e := http.Get(apiURL)
	if e != nil {
		fmt.Println(e)
		return localSyllables(str)

	}
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)
	if e != nil {
		fmt.Println(e)
		return localSyllables(str)
	}
	var count []struct {
		Syllables int `json:"numSyllables"`
	}
	if e := json.Unmarshal(body, &count); e != nil {
		fmt.Println(e)
		return localSyllables(str)
	}
	if len(count) <= 0 {
		fmt.Println(count)
		return localSyllables(str)
	}
	fmt.Println(count, "from api")
	return count[0].Syllables
}

func intermingleRhymes(words []string) [][]string {
	uniqueStrings := make(map[string]struct{})
	var result []string

	for _, str := range words {
		// Check if the string is already in the map
		if _, found := uniqueStrings[str]; !found {
			// If not found, add it to the result and map
			result = append(result, str)
			uniqueStrings[str] = struct{}{}
		}
	}
	var output [][]string
	for i := 0; i < len(result)-2; i++ {
		var temp []string
		temp = append(temp, result[i])
		for j := i + 1; j < len(result)-1; j++ {
			b := doesRhyme(result[i], result[j])
			if b {
				temp = append(temp, result[j])
			}
		}
		output = append(output, temp)
		temp = nil
	}
	return output
}
func rhymeGroupsToString(t [][]string) string {
	var result string
	for _, innerArray := range t {
		for _, str := range innerArray {
			result += " " + str
		}
		result += "\n"
	}
	return result
}

func splitString(input string) []Stanza {
	line := ""
	var stanza Stanza
	var output []Stanza
	for i := 0; i < len(input); i++ {
		teststring := input[i : i+1]
		if teststring != "\n" {
			line += teststring
			if i == len(input)-1 {
				stanza.lines = append(stanza.lines, line)
				stanza.endWords = setEndWords(stanza)
				output = append(output, stanza)
				break
			}

			continue
		}
		stanza.lines = append(stanza.lines, line)
		line = ""
		i++
		if i < len(input) && input[i:i+1] == "\n" {
			stanza.endWords = setEndWords(stanza)
			output = append(output, stanza)
			stanza.lines = nil
			stanza.endWords = nil
			continue
		}
		i--
	}
	return output
}
func removePunctuationFromWord(word string) string {
	removePunctuationFromWordEnd := func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}

	return strings.Map(removePunctuationFromWordEnd, word)
}
func setEndWords(stanza Stanza) []string {
	var output []string
	for i := len(stanza.lines) - 1; i >= 0; i-- {
		temp := strings.Split(stanza.lines[i], " ")

		result := removePunctuationFromWord(temp[len(temp)-1])

		output = append(output, result)
	}
	return output
}
func getAllEndWords(stanzas []Stanza) []string { // formatting is fort testing purposes
	var output []string
	for _, stanza := range stanzas {

		output = append(output, stanza.endWords...)

	}
	return output
}
func stanzasTostrings(stanzas []Stanza) string { // formatting is fort testing purposes
	var output string
	for _, stanza := range stanzas {
		for _, line := range stanza.lines {
			output += line + "\n"
		}
		output += ",\n"
	}
	return output
}
