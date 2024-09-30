package config

import (
	"flag"
)

func init() {
	db := flag.Bool("db", false, "")

	flag.Parse()

	// connect
	loadEnv()
	connectPostgresql(*db)
	initSmptAuth()
}
