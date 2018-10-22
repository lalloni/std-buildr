package cmd

import (
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqleve"
)

var createEventual = &cobra.Command{
	Use:     "create-eventual",
	Short:   "Create standar eventual workspace with the appropriate structure",
	PreRunE: chain(parseScripts),
	RunE:    chain(runCreateEventual),
}

var (
	dmls *[]string
	ddls *[]string
	dcls *[]string
)

func init() {

	initConfig()
	createEventual.Flags().StringP("issue-id", "T", "", "Issue id. Only for SQL Eventual Applications")
	must(viper.BindPFlag("buildr.issue-id", createEventual.Flags().Lookup("issue-id")))

	dmls = createEventual.Flags().StringArray("dml", []string{}, "dml files to be created")
	ddls = createEventual.Flags().StringArray("ddl", []string{}, "ddl files to be created")
	dcls = createEventual.Flags().StringArray("dcl", []string{}, "dcl files to be created")

	rootCmd.AddCommand(createEventual)

}

func index(s string, ss []string) int {
	for i, a := range ss {
		if a == s {
			return i
		}
	}
	return -1
}

func scripts(ss []string, t string) []sqleve.Script {
	a := []sqleve.Script{}
	for _, s := range ss {
		a = append(a, sqleve.Script{Name: s, Type: t})
	}
	return a
}

func parseScripts(ctx *context.Context) error {

	loadProjectConfig(nil, nil)

	ss := append(scripts(*dmls, "dml"), append(scripts(*dcls, "dcl"), scripts(*ddls, "ddl")...)...)

	sort.Slice(ss, func(a, b int) bool {
		return index(ss[a].Name, os.Args) < index(ss[b].Name, os.Args)
	})

	viper.Set("buildr.scripts", ss)

	return nil

}

func runCreateEventual(ctx *context.Context) error {

	cfg := viper.Get("buildr.config").(*config.Config)

	err := initializer.CreateEventual(cfg)
	if err != nil {
		return errors.Wrap(err, "creating eventual worckspace")
	}

	return nil

}
