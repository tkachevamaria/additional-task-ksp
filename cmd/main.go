package main

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"additional-task-ksp/internal/tests"
)

func main() {
	// БД
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// сервис и хэндлер
	service := tests.NewService(db)
	handler := tests.NewHandler(service)

	// роутер
	router := gin.Default()
	router.Use(cors.Default())

	//
	router.GET("/tests", handler.GetAllTests)
	router.GET("/tests/:id", handler.GetTestByID)
	router.POST("/tests/:id/submit", handler.SubmitTest)

	// запуск
	router.Run(":8080")
}
