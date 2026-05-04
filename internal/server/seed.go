package server

import (
	"database/sql"
	"log"
)

func SeedDatabase(db *sql.DB) error {
	log.Println("🔍 [SeedDatabase] Проверяю, есть ли данные в БД...")

	// Проверяем, есть ли уже результаты (а не тесты!)
	var resultsCount int
	err := db.QueryRow("SELECT COUNT(*) FROM results").Scan(&resultsCount)
	if err != nil {
		log.Printf("[SeedDatabase] Ошибка проверки результатов: %v", err)
		return err
	}

	log.Printf("[SeedDatabase] В таблице results найдено %d записей", resultsCount)

	// Если результаты уже есть - пропускаем всё
	if resultsCount >= 12 {
		log.Println("[SeedDatabase] Все 12 результатов уже существуют, пропускаем заполнение")
		return nil
	}

	// Проверяем тесты отдельно
	var testsCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tests").Scan(&testsCount)
	if err != nil {
		log.Printf("[SeedDatabase] Ошибка проверки тестов: %v", err)
		return err
	}

	log.Printf("[SeedDatabase] В таблице tests найдено %d записей", testsCount)

	// Вставка теста, если его нет
	if testsCount == 0 {
		log.Println("[SeedDatabase] Добавляю тест...")
		_, err = db.Exec("INSERT INTO tests (id, title) VALUES (?, ?)", 1, "Волшебный тест личности")
		if err != nil {
			log.Printf("[SeedDatabase] Ошибка вставки теста: %v", err)
			return err
		}
		log.Println("[SeedDatabase] Тест добавлен")
	} else {
		log.Println("[SeedDatabase] Тест уже существует, пропускаю")
	}

	// Вставка вопросов (проверяем отдельно)
	var questionsCount int
	err = db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&questionsCount)
	if err != nil {
		log.Printf("[SeedDatabase] Ошибка проверки вопросов: %v", err)
		return err
	}

	if questionsCount == 0 {
		log.Println("[SeedDatabase] Добавляю вопросы...")
		questions := []struct {
			id     int
			testID int
			text   string
		}{
			{1, 1, "Верите ли вы в Знаки Зодиака?"},
			{2, 1, "Как вы выбираете себе пароль?"},
			{3, 1, "Почему…"},
			{4, 1, "Вы верите в Бога?"},
			{5, 1, "В СССР всё делали…"},
			{6, 1, "Переведите по номеру телефона +79518562061 сумма -  решение примера:"},
		}

		for _, q := range questions {
			_, err := db.Exec("INSERT INTO questions (id, test_id, text) VALUES (?, ?, ?)",
				q.id, q.testID, q.text)
			if err != nil {
				log.Printf("[SeedDatabase] Ошибка вставки вопроса id=%d: %v", q.id, err)
				return err
			}
		}
		log.Printf("[SeedDatabase] Добавлено %d вопросов", len(questions))
	} else {
		log.Printf("[SeedDatabase] Вопросы уже существуют (%d), пропускаю", questionsCount)
	}

	// Вставка ответов
	var answersCount int
	err = db.QueryRow("SELECT COUNT(*) FROM answers").Scan(&answersCount)
	if err != nil {
		log.Printf("[SeedDatabase] Ошибка проверки ответов: %v", err)
		return err
	}

	if answersCount == 0 {
		log.Println("[SeedDatabase] Добавляю ответы...")
		answers := []struct {
			id         int
			questionID int
			text       string
		}{
			{1, 1, "Да"},
			{2, 1, "Нет"},
			{3, 1, "Только по понедельникам"},
			{4, 2, "Да"},
			{5, 2, "Не"},
			{6, 2, "{пароль пользователя}"},
			{7, 3, "Разговор со мной начинается не с поклона?"},
			{8, 3, "Зачем..."},
			{9, 3, "Я здесь?"},
			{10, 4, "Да"},
			{11, 4, "Нет"},
			{12, 4, "Пусть он верит в меня :3"},
			{13, 5, "Не спеша"},
			{14, 5, "Не дыша"},
			{15, 5, "Черемша"},
			{16, 5, "Четыре карандаша"},
			{17, 5, "С лицом вождя, с душой моржа"},
			{18, 6, "Да"},
			{19, 6, "Нет"},
			{20, 6, "Я не знаю"},
		}

		for _, a := range answers {
			_, err := db.Exec("INSERT INTO answers (id, question_id, text) VALUES (?, ?, ?)",
				a.id, a.questionID, a.text)
			if err != nil {
				log.Printf("[SeedDatabase] Ошибка вставки ответа id=%d: %v", a.id, err)
				return err
			}
		}
		log.Printf("[SeedDatabase] Добавлено %d ответов", len(answers))
	} else {
		log.Printf("[SeedDatabase] Ответы уже существуют (%d), пропускаю", answersCount)
	}

	// Вставка результатов
	if resultsCount == 0 {
		log.Println("[SeedDatabase] Добавляю результаты...")
		results := []struct {
			id          int
			testID      int
			title       string
			description string
		}{
			{1, 1, "Овен", "Овен — огненный знак, полный энергии и решительности. Ты прирождённый лидер, который не боится трудностей и всегда идёт вперёд!"},
			{2, 1, "Телец", "Телец — знак стабильности и надёжности. Ты ценишь комфорт, красоту и умеешь наслаждаться жизнью во всех её проявлениях."},
			{3, 1, "Близнецы", "Близнецы — интеллектуалы и мастера общения. Твой ум быстр, а любопытство не знает границ. Ты легко адаптируешься к любым ситуациям!"},
			{4, 1, "Рак", "Рак — знак глубоких чувств и интуиции. Ты обладаешь невероятной эмпатией и способностью создавать уют вокруг себя."},
			{5, 1, "Лев", "Лев — царственный знак, излучающий уверенность и харизму. Ты рождён блистать и вдохновлять окружающих своим примером!"},
			{6, 1, "Дева", "Дева — знак порядка и точности. Твой аналитический ум и внимание к деталям помогают тебе достигать совершенства во всём."},
			{7, 1, "Весы", "Весы — дипломаты и эстеты. Ты стремишься к гармонии, справедливости и умеешь находить баланс даже в самых сложных ситуациях."},
			{8, 1, "Скорпион", "Скорпион — знак страсти и трансформации. Твоя внутренняя сила и решительность способны преодолеть любые препятствия."},
			{9, 1, "Стрелец", "Стрелец — философ и путешественник. Твой оптимизм и жажда приключений вдохновляют всех вокруг на великие свершения!"},
			{10, 1, "Козерог", "Козерог — знак амбиций и дисциплины. Твоё упорство и целеустремлённость помогают тебе достигать невероятных высот."},
			{11, 1, "Водолей", "Водолей — новатор и мечтатель. Твой оригинальный взгляд на мир и стремление к свободе делают тебя уникальной личностью."},
			{12, 1, "Рыбы", "Рыбы — знак интуиции и творчества. Твоя чувствительность и богатое воображение позволяют тебе видеть мир в особенных красках."},
		}

		for _, r := range results {
			log.Printf("  [SeedDatabase] Вставляю результат id=%d: %s", r.id, r.title)
			_, err := db.Exec(
				"INSERT INTO results (id, test_id, title, description) VALUES (?, ?, ?, ?)",
				r.id, r.testID, r.title, r.description,
			)
			if err != nil {
				log.Printf("[SeedDatabase] Ошибка вставки результата id=%d: %v", r.id, err)
				return err
			}
		}
		log.Printf("[SeedDatabase] Добавлено %d результатов", len(results))
	} else {
		log.Printf("[SeedDatabase] Результаты уже существуют (%d), но их меньше 12!", resultsCount)
		log.Println("[SeedDatabase] Нужно удалить старые результаты и вставить новые")
		// Удаляем старые и вставляем заново
		_, err = db.Exec("DELETE FROM results")
		if err != nil {
			log.Printf("[SeedDatabase] Ошибка удаления старых результатов: %v", err)
			return err
		}
		log.Println("[SeedDatabase] Перезапускаю SeedDatabase...")
		return SeedDatabase(db) // Рекурсивно вызываем заново
	}

	log.Printf("[SeedDatabase] База данных успешно заполнена!")
	log.Printf("[SeedDatabase] Итого: тестов=%d, вопросов=%d, ответов=%d, результатов=%d",
		testsCount, questionsCount, answersCount, resultsCount)

	return nil
}
