package main

import (
	"flag"
)

var host, password, dumpfile string
var port, db int

func init() {
	flag.StringVar(&host, "host", "localhost", "Redis host")
	flag.IntVar(&port, "port", 6379, "Redis port")
	flag.IntVar(&db, "db", 0, "Redis database number")
	flag.StringVar(&password, "password", "", "Redis password")
	flag.StringVar(&dumpfile, "dumpfile", "redis-dump.json", "Dump output file")

	flag.Parse()
}
