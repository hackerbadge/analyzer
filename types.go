package main

import "regexp"

type Payload struct {
	Commits []Commit `json:"commits"`
	Sender  struct {
		AvatarUrl string `json:"avatar_url"`
	} `json:"sender"`
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
	Source    string  `json:"source"`
	Username  string  `json:"username"`
	AvatarUrl string  `json:"avatar_url"`
	Tag       string  `json:"tag"`
	Xp        float64 `json:"xp"`
}

type Rule struct {
	Paths []string `yaml:"paths"`
	Tag   string   `yaml:"tag"`
	Xp    float64  `yaml:"xp"`

	initialized bool
	regexps     []*regexp.Regexp
}

func (this *Rule) Init() {
	this.initialized = true

	this.regexps = make([]*regexp.Regexp, len(this.Paths))
	for i, path := range this.Paths {
		this.regexps[i] = regexp.MustCompile(path)
	}
}

func (this *Rule) Apply(c *Commit) bool {
	if !this.initialized {
		this.Init()
	}

	for _, fname := range c.Added {
		if this.applyForFilename(fname) {
			return true
		}
	}

	for _, fname := range c.Modified {
		if this.applyForFilename(fname) {
			return true
		}
	}

	return false
}

func (this *Rule) applyForFilename(name string) bool {
	for _, r := range this.regexps {
		if r.MatchString(name) {
			return true
		}
	}
	return false
}
