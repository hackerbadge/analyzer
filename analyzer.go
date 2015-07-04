package main

import "path/filepath"

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}

type LanguageAnalyzer struct {
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
			promotions = append(promotions, Promotion{source, username, lang, 10})
		}
	}

	return promotions, nil
}

type RulesAnalyzer struct {
}

func (a *RulesAnalyzer) Analyze(commits []Commit) ([]Promotion, error) {
	return []Promotion{}, nil
}

func NewLanguageAnalyzer() Analyzer {
	return &LanguageAnalyzer{}
}

func NewRulesAnalyzer(rules []Rule) Analyzer {
	return &RulesAnalyzer{}
}
