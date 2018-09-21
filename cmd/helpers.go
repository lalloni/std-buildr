package cmd

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
)

func chain(fs ...func(*context.Context) error) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		return chain0(fs...)
	}
}

func chain0(fs ...func(*context.Context) error) error {
	ctx := &context.Context{}
	for _, f := range fs {
		n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf("━━━━━━ running %+v task ", n[strings.Index(n, ".run")+4:]))
		for b.Len() < 80 {
			b.WriteRune('━')
		}
		log.Info(b.String())
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
