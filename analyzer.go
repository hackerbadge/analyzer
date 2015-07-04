package main

type Analyzer interface {
	Analyze(data []Commit) ([]Promotion, error)
}
