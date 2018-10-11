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

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/packager"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the current version of the project",
	RunE:  chain(runClean, runPackage),
}

func init() {

	packageCmd.PersistentFlags().BoolP("allow-dirty", "d", false, "Allow build packages that contain modified files in the working directory")
	must(viper.BindPFlag("buildr.allow-dirty", packageCmd.PersistentFlags().Lookup("allow-dirty")))

	packageCmd.PersistentFlags().BoolP("allow-untagged", "u", false, "Allow build packages that contain that contain commits after the last tag")
	must(viper.BindPFlag("buildr.allow-untagged", packageCmd.PersistentFlags().Lookup("allow-untagged")))

	rootCmd.AddCommand(packageCmd)

}

func runPackage(ctx *context.Context) error {

	git.GetStateIn(ctx)

	log.Infof("version is '%s'", ctx.Build.Version)
	log.Infof("untracked files present: %v", ctx.Build.Untracked)
	log.Infof("changed tracked files present: %v", ctx.Build.Changed)
	log.Infof("uncommited staged files present: %v", ctx.Build.Uncommited)

	cfg := viper.Get("buildr.config").(*config.Config)

	pkgr, err := packager.New(cfg)
	if err != nil {
		return errors.Wrap(err, "getting packager")
	}

	err = pkgr.Package(cfg, ctx)
	if err != nil {
		return errors.Wrap(err, "packaging project")
	}

	return nil
}
