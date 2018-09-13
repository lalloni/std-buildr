// Copyright © 2018 Pablo Lalloni <plalloni@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/packages"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the current version of the project",
	RunE:  chain(runClean, runPackage),
}

func init() {
	rootCmd.AddCommand(packageCmd)
}

func runPackage(ctx *context.Context) error {
	gitversion, err := sh.Output("git", "--version")
	if err != nil {
		return errors.Wrapf(err, "determining git version: %s", gitversion)
	}
	log.Info(gitversion)

		return errors.Errorf("tag name must be prefixed with a 'v' character (found '%s')", ctx.Build.Version)
	if err != nil {
		return errors.Wrapf(err, "getting version from git: %s", ctx.Build.Version)
	}
	log.Infof("version is '%s'", ctx.Build.Version)

	s, err := sh.Output("git", "ls-files", "--exclude-standard", "--others")
	if err != nil {
		return errors.Wrapf(err, "listing untracked files from git: %s", s)
	}
	ctx.Build.Untracked = len(s) > 0
	log.Infof("untracked files present: %v", ctx.Build.Untracked)

	err = sh.Run("git", "diff-files", "--quiet")
	if err != nil {
		ctx.Build.Changed = true
	}
	log.Infof("changed tracked files present: %v", ctx.Build.Changed)

	err = sh.Run("git", "diff-index", "--cached", "--quiet", "HEAD")
	if err != nil {
		ctx.Build.Uncommited = true
	}
	log.Infof("uncommited staged files present: %v", ctx.Build.Uncommited)

	cfg := viper.Get("buildr.config").(*config.Config)

	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		err = packages.PackageOracleSQLEvolutional(ctx, cfg)
	case config.TypeOracleSQLEventual:
		err = packages.PackageOracleSQLEventual(ctx, cfg)
	default:
		err = errors.Errorf("type not implemented: '%s'", cfg.Type)
	}
	return err
}
