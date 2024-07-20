package db_test_database

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
	"github.com/ppp3ppj/go-refactoring-workshop/config"
)

type sqliteDatabase struct {
    *sqlx.DB
}

var (
    sqliteDatabaseInstance *sqliteDatabase
    once                   sync.Once
)

func NewSQLiteDatabase(conf *config.Database) *sqliteDatabase {
    once.Do(func() {
        dsn := fmt.Sprintf("%s_test.db", conf.DBName)

        conn, err := sqlx.Connect("sqlite3", dsn)
        if err != nil {
            panic(err)
        }

        log.Printf("Connected to SQLite database for test %s successfully", conf.DBName)

        sqliteDatabaseInstance = &sqliteDatabase{conn}
    })

    return sqliteDatabaseInstance
}

func (db *sqliteDatabase) Connect() *sqlx.DB {
    return sqliteDatabaseInstance.DB
}

// Close closes the database connection
func (db *sqliteDatabase) Close() error {
    if sqliteDatabaseInstance != nil && sqliteDatabaseInstance.DB != nil {
        if err := sqliteDatabaseInstance.DB.Close(); err != nil {
            return err
        }
        log.Println("Closed SQLite database connection")
    }
    return nil
}
