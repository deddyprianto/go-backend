package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        fmt.Printf("Gagal membaca file .env: %v\n", err)
    }
}

// NewConnection mengembalikan koneksi ke database
func NewConnection() (*sql.DB, error) {
    // Format string koneksi: username:password@tcp(host:port)/dbname
    connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        os.Getenv("DB_USERNAME"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    // Koneksi ke database
    db, err := sql.Open("mysql", connStr)
    if err != nil {
        return nil, fmt.Errorf("gagal koneksi ke database: %v", err)
    }

    return db, nil
}

// PingDatabase untuk testing koneksi
func PingDatabase(db *sql.DB) error {
    err := db.Ping()
    if err != nil {
        return fmt.Errorf("gagal test koneksi: %v", err)
    }
    return nil
}