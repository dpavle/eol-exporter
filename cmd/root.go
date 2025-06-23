/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"eol-exporter/internal/exporter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eol-exporter",
	Short: "Prometheus exporter for OS, kernel, and installed software End-of-Life (EOL) information.",
	Long: `eol-exporter is a Prometheus exporter that exposes End-of-Life (EOL) information about your system.
By default, it collects and exports EOL data for your operating system and kernel.

Additional software or products can be included via a simple plugin system, allowing you to monitor EOL status for any software or product.

The EOL data is pulled from the endoflife.date API (https://endoflife.date/docs/api/v1/) and is refreshed every 24 hours.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		err := exporter.StartExporter()
		if err != nil {
			fmt.Println(err)
		}
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file")
	rootCmd.PersistentFlags().String("listen-port", "3020", "port to start HTTP exporter on")
	rootCmd.PersistentFlags().String("listen-address", "0.0.0.0", "address to start HTTP exporter on")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("listen-port"))
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("listen-address"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
