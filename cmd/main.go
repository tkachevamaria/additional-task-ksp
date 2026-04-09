package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// Подключаемся к БД (если файла нет - создастся)
	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицы (если их нет)
	createTables(db)

	// Дальше работаем с БД
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tests").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Количество тестов в базе: %d\n", count)
}

func createTables(db *sql.DB) {
	// Твой SQL код из файла
	sql := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS tests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS questions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        test_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS answers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        result_tag TEXT NOT NULL,
        FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS results (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        test_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        result_tag TEXT NOT NULL,
        FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
    );`

	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("Ошибка создания таблиц:", err)
	}

	fmt.Println("Таблицы созданы (или уже существуют)")
}
