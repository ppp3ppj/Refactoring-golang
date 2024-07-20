package main

import (
	"log"

	"github.com/ppp3ppj/go-refactoring-workshop/config"
	"github.com/ppp3ppj/go-refactoring-workshop/db"
	db_test_database "github.com/ppp3ppj/go-refactoring-workshop/db/test_database"
)

func main() {
    conf := config.ConfigGetting()
    db := db.NewPostgresDatabase(conf.Database)
    defer func() {
        if err := db.Close(); err != nil {
            log.Fatal("Failed to close database connection: %v", err)
        }
    }()
    // for test sqlite db please remove later.
    sqliteTestDB := db_test_database.NewSQLiteDatabase(conf.Database)
    if err := sqliteTestDB.Close(); err != nil {
        log.Fatal(err)
    }
}
