package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

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
	languageAnalyzer = NewLanguageAnalyzer(config.Source, config.DefaultPoints)

	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	}).Methods("GET")
	r.HandleFunc("/commit", CommitHandler).Methods("POST")
	r.HandleFunc("/import", ImportHandler).Methods("POST")

	if port := os.Getenv("VCAP_APP_PORT"); len(port) != 0 {
		if p, e := strconv.Atoi(port); e == nil && p > 0 {
			config.Port = int(p)
		}
	}

	http.Handle("/", r)
	log.Printf("Listening on port %d\n", config.Port)

	listen := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Println(http.ListenAndServe(listen, nil))
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

	if len(promos) == 0 {
		return
	}

	log.Printf("message for user %s, given points: %f\n", promos[0].Username, promos[0].Points)

	// HACK
	for i := range promos {
		promos[i].AvatarUrl = p.Sender.AvatarUrl
	}

	resp, err := sendToCollector(promos)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	w.Write(respData)
}

// ImportHandler takes github repo name, e.g. hackerbadge/analyzer, and import all commits, push to Collector API
func ImportHandler(w http.ResponseWriter, r *http.Request) {
	var (
		clientId      = "3a758ff9868a3541c9cf"
		clientSecret  = "dc7e30f04713519c02f8730808d10f462163e528"
		queries       = r.URL.Query()
		name          = queries["name"][0]
		singleCommits []GithubSingleCommit

		wg  sync.WaitGroup
		max = 20
		i   = 0
	)

	commitUrls, err := fetchAllCommitURLs(name, clientId, clientSecret)
	if err != nil {
		panic(err)
	}

	// loop and fetch all single commits, collect changed files
	for {
		if i >= len(commitUrls) {
			break
		}

		ch := make(chan GithubSingleCommit, max)
		for j := 0; j < max; j++ {
			if i >= len(commitUrls) {
				break
			}
			wg.Add(1)
			go fetchCommitURL(commitUrls[i], clientId, clientSecret, ch, &wg)
			i++
		}
		wg.Wait()
		close(ch)

		for m := range ch {
			singleCommits = append(singleCommits, m)
		}
	}

	// Send singleCommits to analyzer
	analyzer := NewLanguageAnalyzer(config.Source, config.DefaultPoints)
	promos, err := analyzer.AnalyzeFull(singleCommits)
	if err != nil {
		panic(err)
	}

	if len(promos) == 0 {
		return
	}

	resp, err := sendToCollector(promos)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	log.Println("Response from Collector API:" + string(respData))
	w.Write(respData)
}

func fetchCommitURL(url, clientId, clientSecret string, ch chan GithubSingleCommit, wg *sync.WaitGroup) {
	defer wg.Done()
	url = fmt.Sprintf("%s?client_id=%s&client_secret=%s", url, clientId, clientSecret)
	log.Printf("[DEBUG] Fetching single Commit URL %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	singleCommit := &GithubSingleCommit{}

	// Decoding json response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(singleCommit)
	if err != nil {
		panic(err)
	}

	log.Printf("[DEBUG] Fetched commit %+v\n", *singleCommit)
	ch <- *singleCommit
}

func fetchAllCommitURLs(name, clientId, clientSecret string) ([]string, error) {
	var (
		commitUrls []string
		page       = 1
		perPage    = 50
		err        error
	)

	// loop and fetch all pages of /commits API, collect all URLs of single commits
	for {
		apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/commits?page=%d&per_page=%d&client_id=%s&client_secret=%s", name, page, perPage, clientId, clientSecret)
		log.Printf("[DEBUG] Fetching Commits List from %s\n", apiUrl)

		resp, err := http.Get(apiUrl)
		if err != nil {
			return commitUrls, err
		}
		defer resp.Body.Close()

		// Decoding json response
		decoder := json.NewDecoder(resp.Body)
		githubCommits := []GithubCommit{}
		err = decoder.Decode(&githubCommits)
		if err != nil {
			return commitUrls, err
		}

		for _, githubCommit := range githubCommits {
			commitUrls = append(commitUrls, githubCommit.Url)
		}

		// Stop fetching if there is no more commits
		// TODO remove break here
		break
		if len(githubCommits) == 0 {
			break
		}

		page++
	}

	return commitUrls, err
}

func sendToCollector(promos []Promotion) (resp *http.Response, err error) {
	log.Printf("Sending %d promotions to %s\n", len(promos), config.CollectorApi)
	data, err := json.Marshal(promos)
	if err != nil {
		return nil, err
	}
	log.Printf("Posting to Collector API %s, payload=%s\n", config.CollectorApi, data)
	r := bytes.NewReader(data)
	resp, err = http.Post(config.CollectorApi, "application/json", r)
	fmt.Println("sending to collector finished. Sent promos: %d", len(promos))
	fmt.Println("error: %v", err)
	return
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
