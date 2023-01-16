/*
Package cmd has the commands
*/
package cmd

import (
    "github.com/Eldius/mock-server-go/internal/config"
    "log"
    "os"

    "github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "mock-server-go",
    Short: "A simple mock server engine",
    Long:  `A simple mock server engine.`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    //Run: func(cmd *cobra.Command, args []string) {
    //	server.Start(8080, 8081)
    //},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Here you will define your flags and configuration settings.
    // Cobra supports persistent flags, which, if defined here,
    // will be global for your application.

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mock-server-go.yaml)")

    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    config.Setup(cfgFile)
}
