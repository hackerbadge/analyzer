package main

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}

func NewLanguageAnalyzer() Analyzer {
	return nil
}

func NewRulesAnalyzer(rules []Rule) Analyzer {
	return nil
}
