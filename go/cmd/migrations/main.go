package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
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

var (
	secretCache, _ = secretcache.New()
)

func handleRequest(ctx context.Context, event json.RawMessage) {
	// Decode the JSON event into a struct
	var eventData struct {
		Args []string `json:"args"`
	}
	if err := json.Unmarshal(event, &eventData); err != nil {
		exitf("error unmarshaling event: %v", err)
	}
	userName := os.Getenv("POSTGRES_USER")
	host := os.Getenv("POSTGRES_HOST")
	fmt.Printf("database host: %v\n", host)
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	database := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")
	secretArn := os.Getenv("DB_MASTER_SECRET_ARN")
	if secretArn != "" {
		fmt.Printf("getting password from secret: %v\n", secretArn)
		result, err := secretCache.GetSecretString(secretArn)
		if err != nil {
			fmt.Printf("error retrieving database secret: %v\n", err)
			exitf("error retrieving database secret: %v", err)
		}
		var passwordData struct {
			Password string `json:"password"`
			Username string `json:"username"`
		}
		if err := json.Unmarshal(event, &passwordData); err != nil {
			exitf("error unmarshaling event: %v", err)
		}
		fmt.Printf("password from secret: %v\n", result)
		password = passwordData.Password
	}

	fmt.Println(userName, password, host, dbPort, database)

	db := pg.Connect(&pg.Options{
		User:     userName,
		Database: database,
		Password: password,
		Addr:     fmt.Sprintf("%s:%s", host, dbPort),
		OnConnect: func(ctx context.Context, cn *pg.Conn) error {
			fmt.Println("connected to pg database")
			return nil
		},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	err := db.Ping(ctx)
	if err != nil {
		fmt.Println("error pinging database")
	}
	if len(eventData.Args) == 0 || eventData.Args[0] != "init" {
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
	oldVersion, newVersion, err := migrations.Run(db, eventData.Args...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func main() {
	fmt.Printf("Executing migration...\n")
	lambda_mode := os.Getenv("LAMBDA_MODE")
	if lambda_mode == "" {
		fmt.Printf("Executing in lambda mode\n")
		ctx := context.Background()
		flag.Usage = usage
		flag.Parse()
		// Create a JSON message with args
		eventData := map[string]interface{}{
			"args": flag.Args(),
		}
		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			log.Fatal("Failed to marshal event JSON:", err)
		}

		handleRequest(ctx, json.RawMessage(eventJSON))
	} else {
		lambda.Start(handleRequest)
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
