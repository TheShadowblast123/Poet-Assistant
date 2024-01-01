package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
)

type JsonReturn int

const (
	rhymes JsonReturn = iota
	syllables
	synonyms
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

const rhymeUrl string = `https://api.datamuse.com/words?rel_rhy=`
const nearRhymeUrl string = `https://api.datamuse.com/words?rel_nry=`
const baseDatamuseUrl string = `https://api.datamuse.com/words`
const synonymsUrl string = `https://api.datamuse.com/words?ml=`

func getSynonyms(str string) ([]string, []string) {
	result := validateJsonApiRequest(synonymsUrl+str, rhymes)
	if !result.IsJson {
		fmt.Println(result.e)
		return nil, nil
	}
	var smallOut []string
	var bigOut []string
	for i, r := range result.result {
		bigOut = append(bigOut, r.Word)
		if i < 5 {
			smallOut = append(smallOut, r.Word)
		}
	}
	return smallOut, bigOut
}
func syllableCount(str string) int {
	apiURL := fmt.Sprintf("%s?sp=%s&md=s", baseDatamuseUrl, str)

	res, e := http.Get(apiURL)
	if e != nil {
		fmt.Println(e)
		return -1

	}
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)
	if e != nil {
		fmt.Println(e)
		return -1
	}
	var count []struct {
		Syllables int `json:"numSyllables"`
	}
	if e := json.Unmarshal(body, &count); e != nil {
		fmt.Println(e)
		return -1
	}
	if len(count) <= 0 {
		return -200
	}
	return count[0].Syllables
}

func validateJsonApiRequest(url string, structure JsonReturn) JsonOrError {
	res, e := http.Get(url)
	if e != nil {

		return JsonOrError{false, nil, e}
	}
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)
	if e != nil {
		return JsonOrError{false, nil, e}
	}

	var words []struct {
		Word string `json:"word"`
	}

	switch structure {
	case rhymes:
		if e := json.Unmarshal(body, &words); e != nil {
			return JsonOrError{false, nil, e}
		}

		return JsonOrError{true, words, nil}

	default:
		return JsonOrError{false, nil, nil}
	}
}
func doesRhyme(a string, b string) bool {

	wordA := removePunctuationFromWord(a)
	wordB := removePunctuationFromWord(b)
	if wordA == wordB {
		return true
	}
	resulta := validateJsonApiRequest(rhymeUrl+wordA, rhymes)
	if !resulta.IsJson {
		resultb := validateJsonApiRequest(rhymeUrl+wordB, rhymes)
		if !resultb.IsJson {
			fmt.Print(resulta.e, resultb.e)
			return false
		}
		for _, w := range resultb.result {
			if w.Word == wordA {
				return true
			}
		}
		fmt.Print(resulta.e)
		return false
	}
	for _, w := range resulta.result {
		if w.Word == wordB {
			return true
		}
	}
	resultb := validateJsonApiRequest(rhymeUrl+wordB, rhymes)
	if !resultb.IsJson {
		fmt.Println(resultb.e)
		return false
	}
	for _, w := range resultb.result {
		if w.Word == wordB {
			return true
		}
	}
	resulta = validateJsonApiRequest(nearRhymeUrl+wordA, rhymes)
	if !resulta.IsJson {
		resultb = validateJsonApiRequest(nearRhymeUrl+wordB, rhymes)
		if !resultb.IsJson {
			fmt.Print(resulta.e, resultb.e)
			return false
		}
		for _, w := range resultb.result {
			if w.Word == wordA {
				return true
			}
		}
		fmt.Print(resulta.e)
		return false
	}
	for _, w := range resulta.result {
		if w.Word == wordB {
			return true
		}
	}
	resultb = validateJsonApiRequest(nearRhymeUrl+wordB, rhymes)
	if !resultb.IsJson {
		fmt.Println(resultb.e)
		return false
	}
	for _, w := range resultb.result {
		if w.Word == wordB {
			return true
		}
	}
	return false

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
func getRhymes(str string) ([]string, []string) {
	var smallOut []string
	var bigOut []string
	result := validateJsonApiRequest(rhymeUrl+str, rhymes)
	if result.IsJson {
		if len(result.result) > 9 {
			for i, w := range result.result {
				bigOut = append(bigOut, w.Word)
				if i < 5 {
					smallOut = append(smallOut, w.Word)
				}
			}
			return smallOut, bigOut
		} else if len(result.result) > 5 {
			for i, w := range result.result {
				bigOut = append(bigOut, w.Word)
				if i < 5 {
					smallOut = append(smallOut, w.Word)
				}
			}

			result = validateJsonApiRequest(nearRhymeUrl+str, rhymes)
			if result.IsJson {
				for _, w := range result.result {
					bigOut = append(bigOut, w.Word)
				}
				return smallOut, bigOut
			}

		} else {
			for _, w := range result.result {
				temp := w.Word
				bigOut = append(bigOut, temp)
				smallOut = append(smallOut, temp)

			}
			result = validateJsonApiRequest(nearRhymeUrl+str, rhymes)
			if result.IsJson {
				for i, w := range result.result {
					bigOut = append(bigOut, w.Word)
					if i < 5 {
						smallOut = append(smallOut, w.Word)
					}
				}
				return smallOut, bigOut
			}

		}
	}
	result = validateJsonApiRequest(nearRhymeUrl+str, rhymes)
	if !result.IsJson {
		return nil, nil
	}
	if len(result.result) < 5 {
		for _, w := range result.result {
			bigOut = append(bigOut, w.Word)
			smallOut = append(smallOut, w.Word)

		}

	} else if len(result.result) < 9 {
		for i, w := range result.result {
			bigOut = append(bigOut, w.Word)
			if i < 5 {
				smallOut = append(smallOut, w.Word)
			}
		}

	} else {
		for _, w := range result.result {
			temp := w.Word
			bigOut = append(bigOut, temp)
			smallOut = append(smallOut, temp)

		}

	}
	return smallOut, bigOut
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
		for _, word := range stanza.endWords {
			output = append(output, word)
		}
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
