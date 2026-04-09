package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Test struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type ErrorResponse struct {
	Message string `json:"message"` // tag json
}

func main() {
	// Подключаем БД
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создаём роутер
	router := gin.Default()

	// ===== endpoint: GET /tests =====
	router.GET("/tests", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title FROM tests")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		defer rows.Close()

		var tests []Test

		for rows.Next() {
			var t Test
			rows.Scan(&t.ID, &t.Title)
			tests = append(tests, t)
		}

		c.JSON(http.StatusOK, tests)
	})

	// Запуск сервера
	router.Run(":8080")
}
