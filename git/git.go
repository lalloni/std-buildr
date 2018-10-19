package git

import (
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

func Init() error {
	return sh.Run("git", "init")
}

func AddRemote(name, remote string) error {
	return sh.Run("git", "remote", "add", name, remote)
}

func Add(s string) error {
	return sh.Run("git", "add", s)
}

func AddAll() error {
	fs, err := ListUntrackedFilesAndChangedFiles()
	if err != nil {
		return errors.Wrapf(err, "listing untracked and changed files")
	}
	for _, f := range fs {
		if err = Add(f); err != nil {
			return errors.Wrapf(err, "adding file '%s' to git index", f)
		}
	}
	return nil
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
