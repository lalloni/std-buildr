package cmd

import "github.com/spf13/cobra"

func chain(fs ...func() error) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		return chain0(fs...)
	}
}

func chain0(fs ...func() error) error {
	for _, f := range fs {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
