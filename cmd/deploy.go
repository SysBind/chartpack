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
	"fmt"
	"github.com/SysBind/chartpack/domain"
	"github.com/SysBind/chartpack/infrastructure"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy images and charts from current pack",
	Long: `Copy over images packed with "pack" command to target machines,
Load them into Docker, and deploy (upgrade or install) the charts`,
	Run: func(cmd *cobra.Command, args []string) {
		var nodes []domain.Node

		// Load Nodes Information
		nodes = infrastructure.Nodes()
		for _, node := range nodes {
			fmt.Printf("%s <%s> (Master: %t)\n", node.Hostname, node.Ip, node.IsMaster)
		}

		// Load Charts from local dir
		loader := infrastructure.LocalLoader{Path: "."}
		charts := loader.Load()

		minions := domain.FilterNodes(nodes, func(node domain.Node) bool { return !node.IsMaster })

		for _, chart := range charts {
			fmt.Printf("chart: %s\n", chart.Name)

			files, err := ioutil.ReadDir(chart.Name)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".tar") {
					infrastructure.CopyToNodes(minions, chart.Name+"/"+file.Name())
					infrastructure.RemoteLoadImage(minions, "/"+file.Name())
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
