package main

import (
	"fmt"
	"log"

	"github.com/ppp3ppj/go-refactoring-workshop/config"
	db_test_database "github.com/ppp3ppj/go-refactoring-workshop/db/test_database"
)

func main() {
    conf := config.ConfigGetting()
    _ = conf
    fmt.Printf("%s %s", conf.AppInfo.Name, conf.Database.Host)
    sqliteTestDB := db_test_database.NewSQLiteDatabase(conf.Database)
    if err := sqliteTestDB.Close(); err != nil {
        log.Fatal(err)
    }
}
