package git

import (
	"fmt"
	"regexp"

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

func Initialize() error {
	return sh.Run("git", "init")
}

func AddRemote(remote string) error {
	return sh.Run("git", "remote", "add", "origin", remote)
}

func Add(s string) error {
	return sh.Run("git", "add", s)
}

func AddAll() ([]string, error) {
	fs, err := ListUntrackedFilesAndChangedFiles()

	if err != nil {
		return fs, errors.Wrapf(err, "listing untracked and changed files")
	}

	for _, element := range fs {
		err = Add(element)
		if err != nil {
			return fs, err
		}
	}

	return fs, nil
}

func Commit(m string) error {
	msg := fmt.Sprintf("'%s'", m)
	return sh.Run("git", "commit", "-m", msg)
}

func ListUntrackedFilesAndChangedFiles() ([]string, error) {
	re := regexp.MustCompile(`\r?\n`)

	s := []string{}

	changed, err := sh.Output("git", "diff-files", "--name-only")
	if err != nil {
		return s, errors.Wrapf(err, "listing changed files from git: %s", changed)
	}

	s = append(s, re.Split(changed, -1)...)

	untrackedFiles, err := sh.Output("git", "ls-files", "--exclude-standard", "--others")
	if err != nil {
		return s, errors.Wrapf(err, "listing untracked files from git: %s", untrackedFiles)
	}

	s = append(s, re.Split(untrackedFiles, -1)...)

	return s, nil
}
