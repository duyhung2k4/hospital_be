package config

import (
	"flag"

	"github.com/go-chi/jwtauth/v5"
)

func init() {
	db := flag.Bool("db", false, "")
	jwt = jwtauth.New("HS256", []byte("hospital"), nil)

	flag.Parse()

	// connect
	loadEnv()
	connectPostgresql(*db)
	initSmptAuth()
}
