# influxcat

A simple tool to dump and restore InfluxDB timeseries. This is **NOT** a replacement for InfluxDB's [backup and restore](https://docs.influxdata.com/influxdb/v1.3/administration/backup_and_restore/) feature as this tool is very naive and was made more to restore data from an InfluxDB instance to another for testing and/or debugging purposes.

## Installing

If you have Go installed on your computer:

`go install github.com/admiralobvious/influxcat`

If you don't have Go installed, download the latest release binary from [here](https://github.com/admiralobvious/influxcat/releases/latest).


## Using

View available commands and help:

`$ influxcat`

```
A tool to dump and restore InfluxDB timeseries

Usage:
  influxcat [command]

Available Commands:
  dump        Dumps a timeseries from a database
  help        Help about any command
  restore     Restores a timeseries into a database
  version     Prints the version number

Flags:
  -a, --addr string       addr should be of the form "http://host:port" (default "http://localhost:8086")
  -c, --config string     config file (default is $HOME/.influxcat.yaml)
  -d, --database string   database is the InfluxDB database, required
  -h, --help              help for influxcat
  -p, --password string   password is the InfluxDB password, optional
  -s, --series string     series is the InfluxDB timeseries, required
  -u, --username string   username is the InfluxDB username, optional

Use "influxcat [command] --help" for more information about a command.
```

Dumping a timeseries to a file:

```
$ influxcat dump -d mydb -s myts data.json
done dumping mydb.myts (1337 rows)
```

Restoring a timeseries from a file (the target database has to exist):

```
$ influxcat restore -d mydb -s myts data.json
finished restoring mydb.myts (1337 rows)
```
