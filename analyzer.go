package main

import "path/filepath"

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}

type LanguageAnalyzer struct {
	xp float64
}

// Analyze commits coming from a single source
func (a *LanguageAnalyzer) Analyze(commits []Commit) ([]Promotion, error) {

	exts := map[string]string{
		".php":  "php",
		".cpp":  "cpp",
		".c":    "c",
		".go":   "golang",
		".py":   "python",
		".rb":   "ruby",
		".css":  "css",
		".html": "html",
		".sh":   "bash",
	}

	var source string = "github"
	var promotions []Promotion
	var langsByUser = make(map[string]map[string]int)
	var lang string

	for _, commit := range commits {
		username := commit.Author.Username
		_, userExists := langsByUser[username]
		if !userExists {
			langsByUser[username] = make(map[string]int)
		}
		for _, file := range commit.Added {
			ext := filepath.Ext(file)
			lang, _ = exts[ext]
			if lang != "" {
				langsByUser[username][lang] = 1
			}
		}
		for _, file := range commit.Modified {
			ext := filepath.Ext(file)
			lang, _ = exts[ext]
			if lang != "" {
				langsByUser[username][lang] = 1
			}
		}
	}

	for username, langs := range langsByUser {
		for lang, _ := range langs {
			promotions = append(promotions, Promotion{source, username, lang, a.xp})
		}
	}

	return promotions, nil
}

func NewLanguageAnalyzer(defaultXp float64) Analyzer {
	return &LanguageAnalyzer{defaultXp}
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
				Username: commit.Author.Username,
				Tag:      rule.Tag,
				Xp:       rule.Xp,
			})

		}
	}
	return promos
}
