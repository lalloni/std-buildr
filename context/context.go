package context

type Context struct {
	Build     Build
	Artifacts []Artifact
}

func (c *Context) AddArtifact(file string) {
	c.Artifacts = append(c.Artifacts, Artifact{File: file})
}

type Build struct {
	Version    string
	Prerelease string
	Untracked  bool
	Changed    bool
	Uncommited bool
}

func (p *Build) Dirty() bool {
	return p.Untracked || p.Uncommited || p.Changed
}

func (p *Build) String() string {
	s := p.Version + p.Prerelease
	if p.Dirty() {
		s = s + "-dirty"
	}
	return s
}

type Artifact struct {
	File string
}
