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
	Short: "Create standar project with the appropriate structure according to the type of project indicated",
	RunE:  chain(runCreateProject),
}

func init() {

	createProject.PersistentFlags().StringP("system-id", "s", "", "System Id")
	must(viper.BindPFlag("buildr.system-id", createProject.PersistentFlags().Lookup("system-id")))

	createProject.PersistentFlags().StringP("applicationId", "a", "", "Application ID")
	must(viper.BindPFlag("buildr.application-id", createProject.PersistentFlags().Lookup("applicationId")))

	createProject.PersistentFlags().StringP("type", "t", "", "Project type")
	must(viper.BindPFlag("buildr.type", createProject.PersistentFlags().Lookup("type")))

	createProject.PersistentFlags().StringP("tracker-id", "T", "", "Issue tracker id. Only for SQL Eventual Applications")
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
