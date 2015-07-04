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
	got, err := a.Analyze(mockCommits)

	if err != nil {
		t.Fatalf("Failed with error: %#v", err)
	}

	want := []Promotion{
		Promotion{"github", "solidfoxrock", "python", 10},
		Promotion{"github", "solidfoxrock", "bash", 10},
		Promotion{"github", "solidfoxrock", "golang", 10},
	}

	if len(got) != len(want) {
		t.Fatalf("Want count = %v, got count = %v", len(want), len(got))
	}

	for i, p := range want {
		if got[i].Source != p.Source || got[i].Tag != p.Tag || got[i].Username != p.Username || got[i].Xp != p.Xp {
			t.Fatalf("Item %v, want %#v, got %#v", i, p, got[i])
		}
	}
}
