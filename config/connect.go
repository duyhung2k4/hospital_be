package config

import (
	"app/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgresql(migrate bool) error {
	var err error
	dns := fmt.Sprintf(
		`
			host=%s
			user=%s
			password=%s
			dbname=%s
			port=%s
			sslmode=disable`,
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	dbPsql, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if migrate {
		errMigrate := dbPsql.AutoMigrate(
			&model.Profile{},
			&model.Face{},
			&model.Department{},
			&model.Field{},
			&model.ProfileDepartment{},
			&model.Room{},
			&model.Schedule{},
			&model.Step{},
		)

		if errMigrate != nil {
			return errMigrate
		}
	}

	return err
}
