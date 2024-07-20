package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ppp3ppj/go-refactoring-workshop/config"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
)

type postgresDatabase struct {
    *sqlx.DB
}

var (
    postgresDatabaseInstance *postgresDatabase
    once sync.Once
)

func NewPostgresDatabase(conf *config.Database) *postgresDatabase {
    once.Do(func() {
        dsn := fmt.Sprintf(
            "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
            conf.Host,
            conf.Port,
            conf.User,
            conf.Password,
            conf.DBName,
            conf.SSLMode,
        )

        conn, err := sqlx.Connect("pgx", dsn)
        if err != nil {
            panic(err)
        }

        log.Printf("Connected to postgres database %s successfully", conf.DBName)

        postgresDatabaseInstance = &postgresDatabase{conn}
    })

    return postgresDatabaseInstance
}

func (db *postgresDatabase) Connect() *sqlx.DB {
    return postgresDatabaseInstance.DB
}

// Close closes the database connection
func (db *postgresDatabase) Close() error {
    if postgresDatabaseInstance != nil && postgresDatabaseInstance.DB != nil {
        if err := postgresDatabaseInstance.DB.Close(); err != nil {
            return err
        }
        log.Println("Closed postgres database connection")
    }
    return nil
}
