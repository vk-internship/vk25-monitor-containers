package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS pings (
            id SERIAL PRIMARY KEY,
            ip_address VARCHAR(15) NOT NULL UNIQUE,
            ping_time TIMESTAMP NOT NULL,
            is_success BOOLEAN NOT NULL,
			last_success_time TIMESTAMP
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	addColumnQuery := `
        ALTER TABLE pings
        ADD COLUMN IF NOT EXISTS last_success_time TIMESTAMP;
    `
	_, err = db.Exec(addColumnQuery)
	if err != nil {
		return fmt.Errorf("ошибка добавления столбца last_success_time: %v", err)
	}

	return err
}

func Connect() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения: %v", err)
	}

	log.Println("Успешное подключение к базе данных")

	return db
}
