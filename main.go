package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
var villanelleTemplate = template.Must(template.ParseFiles("templates/villanelle.html"))
var tankaTemplate = template.Must(template.ParseFiles("templates/tanka.html"))
var sestinaTemplate = template.Must(template.ParseFiles("templates/sestina.html"))
var pantoumTemplate = template.Must(template.ParseFiles("templates/pantoum.html"))
var haikuTemplate = template.Must(template.ParseFiles("templates/haiku.html"))
var cinquainTemplate = template.Must(template.ParseFiles("templates/cinquain.html"))
var sonnetTemplate = template.Must(template.ParseFiles("templates/sonnet.html"))

func handlePantoumAdd(w http.ResponseWriter, r *http.Request) {
	i, _ := strconv.Atoi(r.FormValue("stanzas"))
	i++
	temp := fmt.Sprintf(`                <span class = "sub_settings" id = "stanza-count-holder">
	<button hx-get="/add-pantoum" hx-vals='{"stanzas": %d}' hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="addPantoumStanza()">Add Stanza</button>
	<button hx-get="/sub-pantoum" hx-vals='{"stanzas": %d}'  hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="subPantoumStanza()">Remove Stanza</button>
</span>`, i, i)
	w.Write([]byte(temp))

}
func handlePantoumSubtract(w http.ResponseWriter, r *http.Request) {
	i, _ := strconv.Atoi(r.FormValue("stanzas"))
	var temp string
	if i > 2 {
		i--
		if i <= 3 {
			temp = `                <span class = "sub_settings" id = "stanza-count-holder"> 
		<button hx-get="/add-pantoum"  hx-vals='{"stanzas": 2}' hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="addPantoumStanza()">Add Stanza</button>
	</span>`
		} else {
			temp = fmt.Sprintf(`                <span class = "sub_settings" id = "stanza-count-holder"> 
		<button hx-get="/add-pantoum" hx-vals='{"stanzas": %d}' hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="addPantoumStanza()">Add Stanza</button>
		<button hx-get="/sub-pantoum" hx-vals='{"stanzas": %d}'  hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="subPantoumStanza()">Remove Stanza</button>
	</span>`, i, i)
		}
		w.Write([]byte(temp))
	} else {
		i = 2
		temp = `                <span class = "sub_settings" id = "stanza-count-holder"> 
		<button hx-get="/add-pantoum"  hx-vals='{"stanzas": 2}' hx-swap="outerHTML"  hx-target="#stanza-count-holder" onclick="addPantoumStanza()">Add Stanza</button>
	</span>`
		w.Write([]byte(temp))
	}

}
func checkVillaNelleRhymes(w http.ResponseWriter, r *http.Request) {
	output := `<br>`
	endWordsString := r.FormValue("endwords")
	type resultgetter struct {
		Result []string `json:"result"`
	}
	var t resultgetter
	var endWords []string
	var bArr []string
	var aArr []string
	//Nomenclature is specifically for handling villanelle
	//it is otherwise unreasonable
	err := json.Unmarshal([]byte(endWordsString), &t)
	endWords = t.Result
	if err != nil {
		fmt.Println("cannot break it down")
		w.Write([]byte(output))
		fmt.Println("Error:", err)
		return
	}
	aOne := endWords[2]
	aTwo := endWords[4]
	aArr = append(aArr, aOne, aTwo)
	index, _ := strconv.Atoi(r.FormValue("index"))
	index -= 1
	subject := endWords[index]
	if subject == "" {

		w.Write([]byte(output))
		return
	}
	if index != 3 && index < 5 {
		if aOne == "" {

			output = `<div> </p> Please make certain that this rhymes with A1 </p> </br></div>`
			w.Write([]byte(output))
			return
		} else {
			if aTwo == "" {
				output = `<div> </p> Please make certain that this rhymes with A2 </p> </br></div>`
				w.Write([]byte(output))
				return
			}
			if !doesRhyme(aOne, aTwo) && !deepDoesRhyme(aOne, aTwo) {
				if subject == aOne {
					output = fmt.Sprintf(`<div> </p> There, might be a spelling error, or a mistake but please double check this rhymes with %s </p> </br></div>`, aTwo)
					w.Write([]byte(output))
					return
				}
				output = fmt.Sprintf(`<div> </p> There, might be a spelling error, or a mistake but please double check this rhymes with %s </p> </br></div>`, aOne)
				w.Write([]byte(output))
				return
			}
			w.Write([]byte(output))
			return
		}
	}
	isFrombArr := false
	for i := 3; i < 20; i += 3 {
		temp := endWords[i]
		if subject == temp {
			isFrombArr = true
			continue
		}
		if temp == "" {
			continue
		}
		bArr = append(bArr, temp)
	}
	if isFrombArr {
		if len(bArr) == 0 {
			if subject == endWords[3] {
				w.Write([]byte(output))
				return
			}
			output = `<div> </p> Please make certain that this rhymes with the B lines </p> </br></div>`
			w.Write([]byte(output))
			return
		} else {
			for _, b := range bArr {
				if doesRhyme(b, subject) || deepDoesRhyme(b, subject) {
					w.Write([]byte(output))
					return
				}
			}
			output = `<div> </p> Might be an error or a spelling mistake, but please make sure that this rhymes with the B lines </p> </br></div>`
			w.Write([]byte(output))
			return
		}
	} else {
		if endWords[20] != "" {
			if subject != endWords[20] {
				aArr = append(aArr, endWords[20])
			}

		}
		if endWords[19] != "" {
			if subject != endWords[19] {
				aArr = append(aArr, endWords[19])
			}
		}
		for _, a := range aArr {
			if doesRhyme(a, subject) || deepDoesRhyme(a, subject) {
				w.Write([]byte(output))
				return
			}
		}
		for i := 5; i < 18; i += 3 {
			temp := endWords[i]
			if subject == temp {
				continue
			}
			if temp == "" {
				continue
			}
			if doesRhyme(temp, subject) || deepDoesRhyme(temp, subject) {
				w.Write([]byte(output))
				return
			}
		}
		output = "<div> </p> Might be an error or a spelling mistake, but please make sure that this rhymes with the A lines </p> </br></div>"
		w.Write([]byte(output))
		return

	}
}
func setSyllables(w http.ResponseWriter, r *http.Request) {
	temp := r.PostFormValue("line")
	result := syllableCount(removePunctuationFromWord(temp))
	data := strconv.Itoa(result) + " Syllables"
	htmlContent := fmt.Sprintf(`<span class="syllable-count">%s</span>`, data)

	w.Write([]byte(htmlContent))

}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.ExecuteTemplate(w, "index.html", nil)
}
func villanelleHandler(w http.ResponseWriter, r *http.Request) {
	villanelleTemplate.ExecuteTemplate(w, "villanelle.html", nil)
}
func tankaHandler(w http.ResponseWriter, r *http.Request) {
	tankaTemplate.ExecuteTemplate(w, "tanka.html", nil)
}
func pantoumHandler(w http.ResponseWriter, r *http.Request) {
	pantoumTemplate.ExecuteTemplate(w, "pantoum.html", nil)
}
func sestinaHandler(w http.ResponseWriter, r *http.Request) {
	sestinaTemplate.ExecuteTemplate(w, "sestina.html", nil)
}
func haikuHandler(w http.ResponseWriter, r *http.Request) {
	haikuTemplate.ExecuteTemplate(w, "haiku.html", nil)
}
func cinquainHandler(w http.ResponseWriter, r *http.Request) {
	cinquainTemplate.ExecuteTemplate(w, "cinquain.html", nil)
}
func sonnetHandler(w http.ResponseWriter, r *http.Request) {
	sonnetTemplate.ExecuteTemplate(w, "sonnet.html", nil)
}
func main() {
	fmt.Println(loadJSONData())
	fmt.Print("Running")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pantoum", pantoumHandler)
	http.HandleFunc("/haiku", haikuHandler)
	http.HandleFunc("/tanka", tankaHandler)
	http.HandleFunc("/cinquain", cinquainHandler)
	http.HandleFunc("/villanelle", villanelleHandler)
	http.HandleFunc("/sestina", sestinaHandler)
	http.HandleFunc("/sonnet", sonnetHandler)
	http.HandleFunc("/api-syllables", setSyllables)
	http.HandleFunc("/add-pantoum", handlePantoumAdd)
	http.HandleFunc("/sub-pantoum", handlePantoumSubtract)
	http.HandleFunc("/check-rhymes-nelle", checkVillaNelleRhymes)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8000", nil)
}
