package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var InfluxAddr string
var influxUsername string
var influxPassword string
var influxDatabaseName string
var influxSeriesName string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "influxcat",
	Short: "A tool to dump and restore InfluxDB timeseries",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.influxcat.yaml)")
	RootCmd.PersistentFlags().StringVarP(&InfluxAddr, "addr", "a", "http://localhost:8086",
		"addr should be of the form \"http://host:port\"")
	RootCmd.PersistentFlags().StringVarP(&influxUsername, "username", "u", "",
		"username is the InfluxDB username, optional")
	RootCmd.PersistentFlags().StringVarP(&influxPassword, "password", "p", "",
		"password is the InfluxDB password, optional")
	RootCmd.PersistentFlags().StringVarP(&influxDatabaseName, "database", "d", "",
		"database is the InfluxDB database, required")
	RootCmd.PersistentFlags().StringVarP(&influxSeriesName, "series", "s", "",
		"series is the InfluxDB timeseries, required")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".influxcat" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".influxcat")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
