package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/admiralobvious/influxcat/influx"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores a timeseries into a database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You need to specify an input file")
			os.Exit(1)
		}

		conf := influx.InfluxDBConfig{
			Addr:     InfluxAddr,
			Username: influxUsername,
			Password: influxPassword,
		}

		c := influx.NewInfluxDBClient(conf)

		f, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}

		influx.Restore(c, influxDatabaseName, influxSeriesName, f)
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}
