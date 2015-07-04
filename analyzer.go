package main

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}

type LanguageAnalyzer struct {
}

func (a *LanguageAnalyzer) Analyze(commits []Commit) ([]Promotion, error) {
	return []Promotion{}, nil
}

type RulesAnalyzer struct {
}

func (a *RulesAnalyzer) Analyze(commits []Commit) ([]Promotion, error) {
	return []Promotion{}, nil
}

func NewLanguageAnalyzer() Analyzer {
	return &LanguageAnalyzer{}
}

func NewRulesAnalyzer() Analyzer {
	return &RulesAnalyzer{}
}

func NewRulesAnalyzer(rules []Rule) Analyzer {
	return nil
}
