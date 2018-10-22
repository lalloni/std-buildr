package git

import (
	"bufio"
	"bytes"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

func GetStateIn(ctx *context.Context) error {
	var err error
	ctx.Build.Version, err = DescribeVersionInCWD()
	if err != nil {
		return errors.Wrapf(err, "getting version from git: %s", ctx.Build.Version)
	}
	ctx.Build.Untracked, err = UntrackedFilesInCWD()
	if err != nil {
		return errors.Wrap(err, "checking for untracked files")
	}
	ctx.Build.Changed = ChangedFilesInCWD()
	ctx.Build.Uncommited = UncommittedFilesInCWD()
	return nil
}

func UntrackedFilesInCWD() (bool, error) {
	s, err := sh.Output("git", "ls-files", "--exclude-standard", "--others")
	if err != nil {
		return false, errors.Wrapf(err, "listing untracked files from git: %s", s)
	}
	return len(s) > 0, nil
}

func ChangedFilesInCWD() bool {
	return sh.Run("git", "diff-files", "--quiet") != nil
}

func UncommittedFilesInCWD() bool {
	return sh.Run("git", "diff-index", "--cached", "--quiet", "HEAD") != nil
}

func DescribeVersionInCWD() (string, error) {
	s, err := sh.Output("git", "describe", "--abbrev=40", "--tags", "HEAD")
	if err != nil {
		return "", errors.Wrapf(err, "getting git description: %s", s)
	}
	return s, nil
}

func Init() error {
	return sh.Run("git", "init")
}

func AddRemote(name, remote string) error {
	return sh.Run("git", "remote", "add", name, remote)
}

func Add(s string) error {
	return sh.Run("git", "add", s)
}

func CreateOrphanBranch(name string) error {
	return sh.Run("git", "checkout", "--orphan", name)
}

func CreateBranch(name string) error {
	return sh.Run("git", "checkout", "-b", name)
}

func CreateBranchFrom(name string, from string) error {
	return sh.Run("git", "checkout", "-b", name, from)
}

func ExistBranch(name string) (bool, error) {
	branches, err := sh.Output("git", "branch", "-a")
	if err != nil {
		return false, err
	}
	for _, branch := range split(branches) {
		if branch == name {
			return true, nil
		}
	}
	return false, nil
}

func AddAll() error {
	fs, err := ListUntrackedFilesAndChangedFiles()
	if err != nil {
		return errors.Wrapf(err, "listing untracked and changed files")
	}

	for _, element := range fs {

		err = Add(element)
		if err != nil {
			return err
		}
	}
	return nil
}
func Push(remote string, branch string) error {
	return sh.Run("git", "push", remote, branch)
}

func Commit(m string) error {
	return sh.Run("git", "commit", "-m", m)
}

func CommitAddingAll(m string) error {

	if err := AddAll(); err != nil {
		return errors.Wrap(err, "adding files to git index")
	}

	if err := Commit(m); err != nil {
		return errors.Wrap(err, "creating git commit")
	}

	return nil
}

func ListUntrackedFilesAndChangedFiles() ([]string, error) {
	out, err := sh.Output("git", "diff-files", "--name-only")
	if err != nil {
		return nil, errors.Wrapf(err, "listing changed files from git: %s", out)
	}
	s := split(out)
	out, err = sh.Output("git", "ls-files", "--exclude-standard", "--others")
	if err != nil {
		return nil, errors.Wrapf(err, "listing untracked files from git: %s", out)
	}
	s = append(s, split(out)...)
	return s, nil
}

func split(lines string) []string {
	s := bufio.NewScanner(bytes.NewBufferString(lines))
	r := []string{}
	for s.Scan() {
		l := s.Text()
		if len(l) > 0 {
			r = append(r, l)
		}
	}
	return r
}
