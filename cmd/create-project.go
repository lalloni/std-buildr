package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer"
)

var createProject = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new standard-compliant project structure",
	RunE:  chain(runCreateProject),
}

func init() {

	createProject.PersistentFlags().StringP("system-id", "s", "", "System Id")
	must(viper.BindPFlag("buildr.system-id", createProject.PersistentFlags().Lookup("system-id")))

	createProject.PersistentFlags().StringP("application-id", "a", "", "Application ID")
	must(viper.BindPFlag("buildr.application-id", createProject.PersistentFlags().Lookup("application-id")))

	createProject.PersistentFlags().StringP("type", "t", "", "Project type")
	must(viper.BindPFlag("buildr.type", createProject.PersistentFlags().Lookup("type")))

	createProject.PersistentFlags().StringP("tracker-id", "T", "", "Issue tracker ID for Oracle SQL Eventual applications")
	must(viper.BindPFlag("buildr.tracker-id", createProject.PersistentFlags().Lookup("tracker-id")))

	rootCmd.AddCommand(createProject)

}

func runCreateProject(ctx *context.Context) error {

	cfg := &config.Config{}

	init, err := initializer.New(cfg)
	if err != nil {
		return errors.Wrap(err, "getting initializer")
	}

	err = init.Initialize(cfg)
	if err != nil {
		return errors.Wrap(err, "creating project")
	}

	return nil
}
