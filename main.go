package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
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
var config *Config

func main() {
	c, err := NewConfig("etc/app.ini")
	if err != nil {
		panic(err.Error())
	}
	config = c

	rules, err := readRules(config.RulesFile)
	if err != nil {
		panic(err.Error())
	}
	rulesAnalyzer = NewRulesAnalyzer(rules, config.Source)
	languageAnalyzer = NewLanguageAnalyzer(config.Source, config.DefaultXp)

	r := mux.NewRouter()
	r.HandleFunc("/commit", CommitHandler).
		Methods("POST")

	http.Handle("/", r)
}

func readRules(f string) ([]Rule, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	rules := []Rule{}
	err = yaml.Unmarshal(data, &rules)
	return rules, err
}

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	p := Payload{}
	err := decoder.Decode(&p)
	if err != nil {
		panic(err.Error())
	}

	promos, err := Analyze(p.Commits)
	if err != nil {
		panic(err.Error())
	}

	resp, err := sendToCollector(promos)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	w.Write(respData)
}

func sendToCollector(promos []Promotion) (resp *http.Response, err error) {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	if err := enc.Encode(promos); err != nil {
		return nil, err
	}
	return http.Post(config.CollectorApi, "application/json", r)
}

func Analyze(data []Commit) ([]Promotion, error) {
	promotions := []Promotion{}

	languagePromos, err := languageAnalyzer.Analyze(data)
	if err != nil {
		return nil, err
	}

	rulesPromos, err := rulesAnalyzer.Analyze(data)
	if err != nil {
		return nil, err
	}

	promotions = append(languagePromos, rulesPromos...)
	return promotions, nil
}

// AppendUnique appends items to a slice if they do not exist in that slice yet
func AppendUnique(slice []string, elems ...string) (ret []string) {
	ret = slice
	for _, elem := range elems {
		var b bool = true
		for _, s := range slice {
			// fmt.Printf("%+v - %+v\n", s, elem)
			if elem == s {
				b = false
				continue
			}
		}
		if b {
			ret = append(ret, elem)
		}
	}
	return ret
}
