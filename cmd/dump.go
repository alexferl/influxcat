package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/admiralobvious/influxcat/influx"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dumps a timeseries from a database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You need to specify an output path")
			os.Exit(1)
		}

		f, err := os.Create(args[0])
		defer f.Close()

		if err != nil {
			fmt.Printf("Error creating dump file: %v\n", err)
			os.Exit(1)
		}

		conf := influx.InfluxDBConfig{
			Addr:     InfluxAddr,
			Username: influxUsername,
			Password: influxPassword,
		}

		c := influx.NewInfluxDBClient(conf)
		res := influx.Dump(c, influxDatabaseName, influxSeriesName)

		j, _ := json.Marshal(res)
		_, err = f.Write(j)
		if err != nil {
			fmt.Printf("Error writing dump to file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
