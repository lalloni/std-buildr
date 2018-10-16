package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/publisher"
)

// packageCmd represents the package command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the project",
	RunE:  chain(runClean, runPackage, runPublish),
}

func init() {
	initConfig()
	publishCmd.PersistentFlags().StringP("trust", "t", "", "File path with trusted certificate chain (in PEM format)")
	must(viper.BindPFlag("buildr.trust", publishCmd.PersistentFlags().Lookup("trust")))

	rootCmd.AddCommand(publishCmd)
}

func runPublish(ctx *context.Context) error {

	cfg := viper.Get("buildr.config").(*config.Config)

	psh, err := publisher.New(cfg)
	if err != nil {
		return errors.Wrap(err, "getting publisher")
	}

	err = psh.Publish(cfg, ctx)
	if err != nil {
		return errors.Wrap(err, "Publishing project")
	}

	return nil
}
