/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/Eldius/mock-server-go/server"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start mock server",
	Long: `Start mock server.
For example:

mock-server-go start -p 8080 -a 8081

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(startPort, startAdminPort, args[0])
	},
}

var (
	startPort      int
	startAdminPort int
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&startPort, "port", "p", 8080, "Port to use for mock interface -p 8080")
	startCmd.Flags().IntVarP(&startAdminPort, "admin-port", "a", 8081, "Port to use for admin interface -a 8081")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
