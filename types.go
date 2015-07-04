package main

type Payload struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Added     []string `json:"added"`
	Author    User     `json:"author"`
	Committer User     `json:"committer"`
	Distinct  bool     `json:"distinct"`
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
}

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Promotion struct {
	Source   string  `json:"source"`
	Username string  `json:"username"`
	Tag      string  `json:"tag"`
	Xp       float64 `json:"xp"`
}

type Rule struct {
	Dirs  []string `yaml:"dirs"`
	Files []string `yaml:"files"`
	Tag   string   `yaml:"tag"`
}
