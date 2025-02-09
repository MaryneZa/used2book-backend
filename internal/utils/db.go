package utils

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "sync"

    _ "github.com/go-sql-driver/mysql"
)

var (
    db   *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_HOST"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_NAME"),
        )
		// log.Println("Connecting to DB with DSN:", dsn)


        var err error
        db, err = sql.Open("mysql", dsn) // Use sql.Open instead of sqlx.Connect
        if err != nil {
            log.Fatalf("Error connecting to the database: %v", err)
        }

        // Optionally verify the connection is alive
        if err = db.Ping(); err != nil {
            log.Fatalf("Error pinging the database: %v", err)
        }
    })

    return db
}
