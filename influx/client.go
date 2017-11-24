package influx

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

type InfluxDBConfig struct {
	Addr     string
	Username string
	Password string
}

func NewInfluxDBClient(conf InfluxDBConfig) client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.Addr,
		Username: conf.Username,
		Password: conf.Password,
	})

	if err != nil {
		fmt.Printf("Error connection to InfluxDB: %v\n", err)
		os.Exit(1)
	}

	return c
}

// query convenience function to query the database
func query(c client.Client, cmd string, db string) (res []client.Result) {
	q := client.Query{
		Command:  cmd,
		Database: db,
		Chunked:  true,
	}

	defer c.Close()

	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			fmt.Printf("Error query returned an error: %v\n", response.Error())
			os.Exit(1)
		}
		res = response.Results
	} else {
		fmt.Printf("Error querying InfluxDB: %v\n", err)
		os.Exit(1)
	}

	return res
}

type Series struct {
	Name    string          `json:"name"`
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
	Meta    `json:"_meta"`
}

type Meta struct {
	Fields map[int][]string  `json:"fields"`
	Tags   map[string]string `json:"tags"`
}

func Dump(c client.Client, db, series string) []Series {
	validateArgs(db, series)

	cmd := fmt.Sprintf("SELECT * FROM %s", series)
	res := query(c, cmd, db)

	tags := getTags(c, db)
	tagsMap := map[string]string{}
	for _, tag := range tags {
		for idx, column := range res[0].Series[0].Columns {
			if tag == column {
				tagsMap[strconv.Itoa(idx)] = tag
			}
		}
	}

	fields := getFields(c, db)
	fieldsMap := map[int][]string{}
	if len(fields) > 0 {
		for _, field := range fields {
			for idx, column := range res[0].Series[0].Columns {
				if field[0] == column {
					fieldsMap[idx] = field
				}
			}
		}
	}

	var ser []Series

	s := &Series{
		Name:    res[0].Series[0].Name,
		Columns: res[0].Series[0].Columns,
		Values:  res[0].Series[0].Values,
		Meta: Meta{
			Fields: fieldsMap,
			Tags:   tagsMap,
		},
	}

	ser = append(ser, *s)

	return ser
}

func Restore(c client.Client, db, series string, file io.Reader) {
	validateArgs(db, series)

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: db,
	})

	if err != nil {
		fmt.Printf("Error creating new batch points: %v\n", err)
		os.Exit(1)
	}

	dec := json.NewDecoder(file)
	_, err = dec.Token() // read open bracket
	if err != nil {
		fmt.Printf("Error reading opening bracket: %v\n", err)
		os.Exit(1)
	}

	for dec.More() {
		var s Series

		err := dec.Decode(&s)
		if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			os.Exit(1)
		}

		for _, values := range s.Values {
			newFields := map[string]interface{}{}
			newTags := map[string]string{}

			var t time.Time

			for idx := range values {
				t, err = time.Parse(time.RFC3339, values[0].(string))
				if err != nil {
					fmt.Printf("Error formatting date: %v\n", err)
					os.Exit(1)
				}

				if keyInMapSlice(s.Meta.Fields, idx) {
					switch s.Meta.Fields[idx][1] {
					case "float":
						newFields[s.Meta.Fields[idx][0]] = values[idx].(float64)
					case "integer":
						newFields[s.Meta.Fields[idx][0]] = int(values[idx].(float64))
					case "string":
						newFields[s.Meta.Fields[idx][0]] = values[idx].(string)
					case "boolean":
						newFields[s.Meta.Fields[idx][0]] = values[idx].(bool)
					case "timestamp":
						newFields[s.Meta.Fields[idx][0]] = values[idx].(time.Time)
					default:
						fmt.Printf("Error: Unknown data type for value %v\n", values[idx])
						os.Exit(1)
					}

				} else if keyInMap(s.Meta.Tags, strconv.Itoa(idx)) {
					newTags[s.Meta.Tags[strconv.Itoa(idx)]] = values[idx].(string)
				}
			}

			pt, err := client.NewPoint(
				s.Name,
				newTags,
				newFields,
				t,
			)

			if err != nil {
				fmt.Printf("Error creating new point: %v\n", err)
				os.Exit(1)
			}

			bp.AddPoint(pt)
		}

	}

	_, err = dec.Token() // read closing bracket
	if err != nil {
		fmt.Printf("Error reading closing bracket: %v\n", err)
		os.Exit(1)
	}

	if err := c.Write(bp); err != nil {
		fmt.Printf("Error writing data to InfluxDB: %v\n", err)
		os.Exit(1)
	}

}

func validateArgs(db, series string) {
	if db == "" {
		fmt.Println("Error: You need to specify a database with --database name")
		os.Exit(1)
	}

	if series == "" {
		fmt.Println("Error: You need to specify a series with --series name")
		os.Exit(1)
	}
}

func getFields(c client.Client, db string) [][]string {
	cmd := fmt.Sprintf("SHOW FIELD KEYS")
	res := query(c, cmd, db)

	validateTimeSeries(res[0].Series)

	var fields [][]string
	for _, field := range res[0].Series[0].Values {
		fields = append(fields, []string{field[0].(string), field[1].(string)})
	}

	return fields
}

func getTags(c client.Client, db string) []string {
	cmd := fmt.Sprintf("SHOW TAG KEYS")
	res := query(c, cmd, db)

	validateTimeSeries(res[0].Series)

	var tags []string
	for _, tag := range res[0].Series[0].Values {
		tags = append(tags, tag[0].(string))
	}

	return tags
}

func validateTimeSeries(series []models.Row) {
	if len(series) <= 0 {
		fmt.Println("Error: There is no data in this timeseries")
		os.Exit(1)
	}
}

func keyInMap(m map[string]string, key string) bool {
	if _, ok := m[key]; ok {
		return true
	}

	return false
}

func keyInMapSlice(m map[int][]string, key int) bool {
	if _, ok := m[key]; ok {
		return true
	}

	return false
}
