package main

import "path/filepath"
import "fmt"

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
	AnalyzeFull(data []GithubSingleCommit) ([]Promotion, error)
}

type languageAnalyzerImpl struct {
	source string
	xp     float64
}

var exts = map[string]string{
	".php":  "php",
	".java": "java",
	".cpp":  "cpp",
	".c":    "c",
	".go":   "golang",
	".py":   "python",
	".rb":   "ruby",
	".js":   "javascript",
	".css":  "css",
	".html": "html",
	".sh":   "bash",
}

// Analyze commits coming from a single source
func (a *languageAnalyzerImpl) Analyze(commits []Commit) ([]Promotion, error) {
	var promotions []Promotion
	var langsByUser = make(map[string][]string)
	var lang string

	for _, commit := range commits {
		username := commit.Author.Username
		_, userExists := langsByUser[username]
		if !userExists {
			langsByUser[username] = []string{}
		}
		for _, file := range commit.Added {
			ext := filepath.Ext(file)
			lang, _ = exts[ext]
			if lang != "" {
				langsByUser[username] = AppendUnique(langsByUser[username], lang)
			}
		}
		for _, file := range commit.Modified {
			ext := filepath.Ext(file)
			lang, _ = exts[ext]
			if lang != "" {
				langsByUser[username] = AppendUnique(langsByUser[username], lang)
			}
		}
	}

	for username, langs := range langsByUser {
		for _, lang := range langs {
			promotions = append(promotions, Promotion{
				Source:   a.source,
				Username: username,
				Tag:      lang,
				Xp:       a.xp})
		}
	}

	return promotions, nil
}

// Analyze commits coming from a full repo import
func (a *languageAnalyzerImpl) AnalyzeFull(commits []GithubSingleCommit) ([]Promotion, error) {

	var (
		promotions  []Promotion
		langsByUser = make(map[string][]string)
		lang, name  string
	)

	fmt.Println("[Language Analyzer - Analyze Full] Start looping through commits...")
	for _, commit := range commits {
		name = commit.Commit.Author.Name
		_, userExists := langsByUser[name]
		if !userExists {
			langsByUser[name] = []string{}
		}

		for _, file := range commit.Files {
			if file.Status == "modified" || file.Status == "added" {
				ext := filepath.Ext(file.FileName)
				lang, _ = exts[ext]
				if lang != "" {
					langsByUser[name] = AppendUnique(langsByUser[name], lang)
				}
			}
		}
	}

	fmt.Printf("%+v\n\n", langsByUser)
	for name, langs := range langsByUser {
		for _, lang := range langs {
			promotions = append(promotions, Promotion{a.source, name, lang, a.xp})
		}
	}

	fmt.Println("PROMOTIONS")
	for _, p := range promotions {
		fmt.Printf("%+v\n", p)
	}

	return promotions, nil
}

func NewLanguageAnalyzer(source string, defaultXp float64) Analyzer {
	return &languageAnalyzerImpl{source, defaultXp}
}

func NewRulesAnalyzer(rules []Rule, source string) Analyzer {
	return &rulesAnalyzerImpl{rules, source}
}

type rulesAnalyzerImpl struct {
	rules  []Rule
	source string
}

func (this *rulesAnalyzerImpl) Analyze(data []Commit) ([]Promotion, error) {
	result := []Promotion{}
	for i := range data {
		commit := &data[i]
		promos := this.analyzeCommit(commit)

		result = append(result, promos...)
	}

	return result, nil
}

func (this *rulesAnalyzerImpl) analyzeCommit(commit *Commit) []Promotion {
	promos := []Promotion{}
	for i := range this.rules {
		rule := &this.rules[i]
		if rule.Apply(commit) {
			promos = append(promos, Promotion{
				Source:   this.source,
				Username: commit.Author.Email,
				Tag:      rule.Tag,
				Xp:       rule.Xp,
			})

		}
	}
	return promos
}

// Analyze commits coming from a full repo import
func (a *rulesAnalyzerImpl) AnalyzeFull(commits []GithubSingleCommit) ([]Promotion, error) {
	var promotions []Promotion
	return promotions, nil
}
