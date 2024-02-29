/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"skjald/v2/display"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var start_period string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "skjald",
	Short: "skjald is a time management tool",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		display.Main(start_period)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&start_period, "start", "s", "work", "Starting period name")
}
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".skjald")


	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig();
	if err != nil {
		fmt.Fprintln(os.Stderr, "Using config file failed:", viper.ConfigFileUsed())
		os.Exit(1)
	}
}
