/*
Copyright © 5780 (2020) Asaf Ohayon <asaf@sysbind.co.il>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/SysBind/chartpack/domain"
	"github.com/SysBind/chartpack/infrastructure"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack src dest",
	Short: "package a directory of charts",
	Long: `Traverse diretory of helm charts, parsing each chart/values.yaml
to extract list of docker images needed for this chart, download them, and pack it all

pack charts_dir out_dir`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]
		log.Println("scanning " + src)
		loader := infrastructure.LocalLoader{Path: src}

		charts := loader.Load()

		exporter := infrastructure.Exporter{Src: src, Dest: dest}

		domain.Package(charts, exporter)

		// copy self binary over to des
		if err := infrastructure.Copy(os.Args[0], dest+"/chartpack"); err != nil {
			panic(err)
		}
		if err := os.Chmod(dest+"/chartpack", 0744); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(packCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
