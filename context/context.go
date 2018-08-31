package context

import (
	"github.com/Masterminds/semver"
)

type Context struct {
	Build     Build
	Artifacts []Artifact
}

func (c *Context) AddArtifact(file string) {
	c.Artifacts = append(c.Artifacts, Artifact{File: file})
}

type Build struct {
	Version    *semver.Version
	Untracked  bool
	Changed    bool
	Uncommited bool
}

func (p *Build) Dirty() bool {
	return p.Untracked || p.Uncommited || p.Changed
}

func (p *Build) String() string {
	s := p.Version.String()
	if p.Dirty() {
		s = s + "-dirty"
	}
	return s
}

type Artifact struct {
	File string
}
