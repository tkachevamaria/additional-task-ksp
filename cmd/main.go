package main

import (
	"database/sql"
	"log"

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

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL UNIQUE,
    birth_date TEXT
);

CREATE TABLE IF NOT EXISTS tests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT
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
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
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

	// Заполняем базу данных тестовыми данными
	if err := server.SeedDatabase(db); err != nil {
		log.Printf("Предупреждение при заполнении БД: %v", err)
	}

	// сервис и хэндлер
	service := server.NewService(db)
	handler := server.NewHandler(service)

	// роутер
	router := gin.Default()
	router.Use(cors.Default())

	// маршруты
	router.GET("/tests", handler.GetAllTests)
	router.GET("/tests/:id", handler.GetTestByID)
	router.POST("/tests/:id/submit", handler.SubmitTest)
	router.POST("/register", handler.Register)
	router.POST("/check-full-match", handler.CheckFullMatch)
	router.POST("/check-email-exists", handler.CheckEmailExists)
	router.POST("/check-password-owner", handler.CheckPasswordOwner)
	router.POST("/check-email-password", handler.CheckEmailAndPassword)

	// запуск
	log.Println("Сервер запущен на http://localhost:8080")
	router.Run(":8080")
}
