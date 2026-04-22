package main

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"additional-task-ksp/internal/server"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// SQL для создания таблиц
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        birth_date TEXT
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

	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	dbPath := "./testdb.db"

	// Инициализируем базу данных
	db, err := InitDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// сервис и хэндлер
	service := server.NewService(db)
	handler := server.NewHandler(service)

	// роутер
	router := gin.Default()
	router.Use(cors.Default())

	//
	router.GET("/tests", handler.GetAllTests)
	router.GET("/tests/:id", handler.GetTestByID)
	router.POST("/tests/:id/submit", handler.SubmitTest)
	router.POST("/register", handler.Register)

	// запуск
	router.Run(":8080")

}
