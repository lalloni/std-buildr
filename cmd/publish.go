// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current version of the project",
	RunE:  chain(runClean, runPackage, runPublish),
}

func init() {
	rootCmd.AddCommand(publishCmd)
}

func runPublish(ctx *context.Context) error {
	for _, artifact := range ctx.Artifacts {
		log.Infof("Publishing %+v", artifact)
	}
	return nil
}
