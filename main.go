package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/**
Analyzer of programming language
Input:
 git diff from one commit
Collected data:
 - user email
 - filename
Rules:
 - defined in config file as regexps

Output:
 - JSON sent to Collector API

*/

var languageAnalyzer Analyzer
var rulesAnalyzer Analyzer

func main() {
	c, err := NewConfig("etc/app.ini")
	if err != nil {
		panic(err.Error())
	}

	rules, err := readRules("ets/rules.yml")
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/commit", CommitHandler).
		Methods("POST")

	http.Handle("/", r)
}

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	p := Payload{}
	err := decoder.Decode(&p)
	if err != nil {
		panic(err.Error())
	}

}

func Analyze(data []Commit) ([]Promotion, error) {
	promotions := []Promotion{}

	languagePromos := NewLanguageAnalyzer().Analyze(data)
	rulesPromos := NewRulesAnalyzer("").Analyze(data)

}
