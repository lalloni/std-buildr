package context

import "regexp"

var untaggedVersionRegex = regexp.MustCompile(`.*-\d+-g[0-9a-f]{40}$`)

type Context struct {
	Build     Build
	Artifacts []Artifact
}

func (c *Context) AddArtifact(file string, path string, isPrerelease bool) {
	c.Artifacts = append(c.Artifacts, Artifact{File: file, Path: path, IsPrerelease: isPrerelease})
}

type Build struct {
	Version    string
	Untracked  bool
	Changed    bool
	Uncommited bool
}

func (p *Build) Dirty() bool {
	return p.Untracked || p.Uncommited || p.Changed
}

func (p *Build) Untagged() bool {
	return !untaggedVersionRegex.MatchString(p.Version)
}

func (p *Build) String() string {
	s := p.Version
	if p.Dirty() {
		s = s + "-dirty"
	}
	return s
}

type Artifact struct {
	File         string
	Path         string
	IsPrerelease bool
}
