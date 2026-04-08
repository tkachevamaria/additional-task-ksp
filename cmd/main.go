package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// Подключаемся к существующему файлу test.db
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Пробуем считать количество тестов
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tests").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Количество тестов в базе: %d\n", count)
}
