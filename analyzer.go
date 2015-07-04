package main

import "path/filepath"

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}

type LanguageAnalyzer struct {
	source string
	xp     float64
}

// Analyze commits coming from a single source
func (a *LanguageAnalyzer) Analyze(commits []Commit) ([]Promotion, error) {

	exts := map[string]string{
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

func NewLanguageAnalyzer(source string, defaultXp float64) Analyzer {
	return &LanguageAnalyzer{source, defaultXp}
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
