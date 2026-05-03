package server

import (
	"database/sql"
	"log"
)

func SeedDatabase(db *sql.DB) error {
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
		// Добавляй новые вопросы сюда
		// {4, 1, "Твой новый вопрос?"},
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
	}{
		// Вопрос 1
		{1, 1, "Читаю книги или смотрю фильмы"},
		{2, 1, "Встречаюсь с друзьями"},
		{3, 1, "Занимаюсь творчеством"},
		{4, 1, "Путешествую и исследую новое"},
		// Вопрос 2
		{5, 2, "Семья и близкие люди"},
		{6, 2, "Карьера и достижения"},
		{7, 2, "Саморазвитие и знания"},
		{8, 2, "Свобода и приключения"},
		// Вопрос 3
		{9, 3, "Утро"},
		{10, 3, "День"},
		{11, 3, "Вечер"},
		{12, 3, "Ночь"},
	}

	for _, a := range answers {
		_, err := db.Exec("INSERT INTO answers (id, question_id, text) VALUES (?, ?, ?)",
			a.id, a.questionID, a.text)
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
	}{
		{
			1, 1,
			"🌟 Романтик",
			"Вы цените уют, гармонию и душевное тепло. Ваша сильная сторона - умение создавать комфортную атмосферу вокруг себя. Люди тянутся к вам за поддержкой и пониманием.",
		},
		{
			2, 1,
			"🎉 Душа компании",
			"Вы - настоящий лидер и центр притяжения для окружающих. Ваша энергия заражает других, а умение находить общий язык с людьми открывает перед вами множество дверей.",
		},
		{
			3, 1,
			"🎨 Творец",
			"У вас уникальное творческое мышление! Вы видите мир нестандартно и способны создавать нечто прекрасное из обычных вещей. Ваша креативность - ваш главный дар.",
		},
		{
			4, 1,
			"⚡ Искатель приключений",
			"Вы не боитесь перемен и всегда готовы к новым вызовам! Ваша смелость и любопытство ведут вас к невероятным открытиям и захватывающим приключениям.",
		},
	}

	for _, r := range results {
		_, err := db.Exec("INSERT INTO results (id, test_id, title, description) VALUES (?, ?, ?, ?)",
			r.id, r.testID, r.title, r.description)
		if err != nil {
			return err
		}
	}

	log.Printf("✅ База данных успешно заполнена!")
	log.Printf("📊 Добавлено: 1 тест, %d вопросов, %d ответов, %d результатов",
		len(questions), len(answers), len(results))

	return nil
}