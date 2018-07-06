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

package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/cmd"
)

var (
	version = "dev"
	commit  string
	date    string
)

func main() {
	if version == "dev" {
		viper.Set("version", "dev")
	} else {
		viper.Set("version", fmt.Sprintf("%s-%s (%s)", version, commit, date))
	}
	if isatty.IsTerminal(os.Stdout.Fd()) {
		log.SetHandler(cli.New(os.Stdout))
	}
	cmd.Execute()
}
