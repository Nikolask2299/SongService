package migration

import (
	"errors"
	"fmt"
	"musicservice/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrations() {
	confPostgres, confMigrat, err := config.InitConfigMigration()
	if err!= nil {
		fmt.Println("Error initial config migration")
        panic(err)
    }

	m, err := migrate.New(
		"file://"+ config.Dir(confMigrat.MigrationsPath), 
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?x-migrations-table=%s&sslmode=disable", confPostgres.User, confPostgres.Password, confPostgres.Host, confPostgres.Port, confPostgres.DBName, confMigrat.MigrationsTable),
	)
	 
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations were applied to the database")
			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied to the database successfully")
}

