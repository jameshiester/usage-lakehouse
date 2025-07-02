package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run *.go <command> [args]
`

func main() {
	flag.Usage = usage
	flag.Parse()
	userName := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	fmt.Println(userName, password, host, dbPort, database)

	db := pg.Connect(&pg.Options{
		User:     userName,
		Database: database,
		Password: password,
		Addr:     fmt.Sprintf("%s:%s", host, dbPort),
	})
	if len(flag.Args()) == 0 || flag.Args()[0] != "init" {
		// Check if gopg_migrations table exists
		type TableExistsResult struct {
			Exists bool `pg:"exists"`
		}
		var result TableExistsResult
		_, err := db.QueryOne(&result, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'gopg_migrations') as exists")
		if err != nil {
			exitf("error checking if gopg_migrations table exists: %v", err)
		}

		if !result.Exists {
			fmt.Println("gopg_migrations table does not exist, initializing migrations")
			_, _, err := migrations.Run(db, "init")
			if err != nil {
				exitf("error initializing migrations: %v", err)
			}
		}
	}
	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
