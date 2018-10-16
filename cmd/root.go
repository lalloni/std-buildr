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
	"fmt"
	"os"

	"github.com/apex/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "buildr",
	SilenceUsage:      true,
	PersistentPreRunE: preRunRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.buildr.yaml)")

	rootCmd.PersistentFlags().StringP("change-directory", "C", "", "Change to specified directory before executing any actions")
	must(viper.BindPFlag("buildr.change-directory", rootCmd.PersistentFlags().Lookup("change-directory")))

}

func preRunRoot(cmd *cobra.Command, args []string) error {
	cdto := viper.GetString("buildr.change-directory")
	if cdto != "" {
		log.Infof("changing current directoty to %s", cdto)
		err := os.Chdir(cdto)
		if err != nil {
			return errors.Wrapf(err, "changing current directory to %s", cdto)
		}
	}
	return nil
}

func initConfig() error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return errors.Wrap(err, "getting home dir")
		}

		// Search config in home directory with name ".buildr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".buildr")

		f, err := os.Open("buildr.yaml")
		if err != nil {
			return errors.Wrap(err, "opening buildr.yaml")
		}
		defer f.Close()
		c, err := config.Read(f)
		if err != nil {
			return errors.Wrapf(err, "reading project configuration")
		}
		viper.Set("buildr.config", c)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	return nil
}
