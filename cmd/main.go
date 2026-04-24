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

// seedDatabase заполняет базу данных тестовыми данными
func seedDatabase(db *sql.DB) error {
	// Проверяем, есть ли уже тесты
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tests").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("База данных уже содержит данные, пропускаем заполнение")
		return nil
	}

	log.Println("Начинаем заполнение базы данных тестовыми данными...")

	// Вставка теста
	_, err = db.Exec("INSERT INTO tests (id, title) VALUES (?, ?)", 1, "Волшебный тест личности")
	if err != nil {
		return err
	}

	// Вставка вопросов
	questions := []struct {
		id     int
		testID int
		text   string
	}{
		{1, 1, "Как вы проводите свободное время?"},
		{2, 1, "Что для вас важнее всего в жизни?"},
		{3, 1, "Какое время суток вам нравится больше всего?"},
	}

	for _, q := range questions {
		_, err := db.Exec("INSERT INTO questions (id, test_id, text) VALUES (?, ?, ?)",
			q.id, q.testID, q.text)
		if err != nil {
			return err
		}
	}

	// Вставка ответов
	answers := []struct {
		id         int
		questionID int
		text       string
		resultTag  string
	}{
		// Вопрос 1
		{1, 1, "Читаю книги или смотрю фильмы", "romantic"},
		{2, 1, "Встречаюсь с друзьями", "social"},
		{3, 1, "Занимаюсь творчеством", "creative"},
		{4, 1, "Путешествую и исследую новое", "adventurer"},
		// Вопрос 2
		{5, 2, "Семья и близкие люди", "romantic"},
		{6, 2, "Карьера и достижения", "social"},
		{7, 2, "Саморазвитие и знания", "creative"},
		{8, 2, "Свобода и приключения", "adventurer"},
		// Вопрос 3
		{9, 3, "Утро", "creative"},
		{10, 3, "День", "social"},
		{11, 3, "Вечер", "romantic"},
		{12, 3, "Ночь", "adventurer"},
	}

	for _, a := range answers {
		_, err := db.Exec("INSERT INTO answers (id, question_id, text, result_tag) VALUES (?, ?, ?, ?)",
			a.id, a.questionID, a.text, a.resultTag)
		if err != nil {
			return err
		}
	}

	// Вставка результатов
	results := []struct {
		id          int
		testID      int
		title       string
		description string
		resultTag   string
	}{
		{
			1, 1,
			"🌟 Романтик",
			"Вы цените уют, гармонию и душевное тепло. Ваша сильная сторона - умение создавать комфортную атмосферу вокруг себя. Люди тянутся к вам за поддержкой и пониманием.",
			"romantic",
		},
		{
			2, 1,
			"🎉 Душа компании",
			"Вы - настоящий лидер и центр притяжения для окружающих. Ваша энергия заражает других, а умение находить общий язык с людьми открывает перед вами множество дверей.",
			"social",
		},
		{
			3, 1,
			"🎨 Творец",
			"У вас уникальное творческое мышление! Вы видите мир нестандартно и способны создавать нечто прекрасное из обычных вещей. Ваша креативность - ваш главный дар.",
			"creative",
		},
		{
			4, 1,
			"⚡ Искатель приключений",
			"Вы не боитесь перемен и всегда готовы к новым вызовам! Ваша смелость и любопытство ведут вас к невероятным открытиям и захватывающим приключениям.",
			"adventurer",
		},
	}

	for _, r := range results {
		_, err := db.Exec("INSERT INTO results (id, test_id, title, description, result_tag) VALUES (?, ?, ?, ?, ?)",
			r.id, r.testID, r.title, r.description, r.resultTag)
		if err != nil {
			return err
		}
	}

	log.Printf("✅ База данных успешно заполнена!")
	log.Printf("📊 Добавлено: 1 тест, %d вопросов, %d ответов, %d результатов",
		len(questions), len(answers), len(results))

	return nil
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
	if err := seedDatabase(db); err != nil {
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

	// запуск
	log.Println("Сервер запущен на http://localhost:8080")
	router.Run(":8080")
}
