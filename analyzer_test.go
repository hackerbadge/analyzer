package main

import "testing"

// mock input
var mockUser = User{"solidfoxrock@gmail.com", "Kien Nguyen", "solidfoxrock"}
var mockCommit = Commit{
	Added:     []string{"bin/setup.py", "bin/say-hello.sh"},
	Author:    mockUser,
	Committer: mockUser,
	Distinct:  true,
	ID:        "aa45b6ee05606d0c62e580bbde433c43ea1136b7",
	Message:   "Fix a Heisenbug",
	Modified:  []string{"etc/app.ini", "etc/rules.yml", "main.go"},
	Removed:   []string{},
	Timestamp: "123",
	URL:       "https://github.com/hackerbadge/hackerbadge/commit/aa45b6ee05606d0c62e580bbde433c43ea1136b7",
}
var mockCommits = []Commit{mockCommit}

func TestLanguageAnalyzer(t *testing.T) {
	a := NewLanguageAnalyzer()
	a.Analyze(mockCommits)
}
