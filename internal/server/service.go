package server

import (
	"additional-task-ksp/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetAllTests возвращает список всех тестов
func (s *Service) GetAllTests() ([]models.Test, error) {
	query := "SELECT id, title FROM tests ORDER BY id"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tests: %w", err)
	}
	defer rows.Close()

	var tests []models.Test
	for rows.Next() {
		var test models.Test
		if err := rows.Scan(&test.ID, &test.Title); err != nil {
			return nil, fmt.Errorf("failed to scan test: %w", err)
		}
		tests = append(tests, test)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tests, nil
}

// GetTestByID возвращает тест с вопросами и ответами
func (s *Service) GetTestByID(testID int) (*models.TestDetail, error) {
	// Проверяем существует ли тест
	var test models.TestDetail
	query := "SELECT id, title FROM tests WHERE id = ?"
	err := s.db.QueryRow(query, testID).Scan(&test.ID, &test.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	// Получаем вопросы и ответы
	questions, err := s.getQuestionsWithAnswers(testID)
	if err != nil {
		return nil, err
	}
	test.Questions = questions

	return &test, nil
}

// getQuestionsWithAnswers возвращает все вопросы теста с ответами
func (s *Service) getQuestionsWithAnswers(testID int) ([]models.QuestionDetail, error) {
	query := `
        SELECT id, text 
        FROM questions 
        WHERE test_id = ? 
        ORDER BY id
    `
	rows, err := s.db.Query(query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()

	var questions []models.QuestionDetail
	for rows.Next() {
		var q models.QuestionDetail
		if err := rows.Scan(&q.ID, &q.Text); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		// Получаем ответы для этого вопроса
		answers, err := s.getAnswersByQuestionID(q.ID)
		if err != nil {
			return nil, err
		}
		q.Answers = answers
		questions = append(questions, q)
	}

	return questions, nil
}

// getAnswersByQuestionID возвращает все ответы для конкретного вопроса
func (s *Service) getAnswersByQuestionID(questionID int) ([]models.Answer, error) {
	query := `
        SELECT id, text, result_tag 
        FROM answers 
        WHERE question_id = ? 
        ORDER BY id
    `
	rows, err := s.db.Query(query, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers: %w", err)
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		if err := rows.Scan(&answer.ID, &answer.Text, &answer.ResultTag); err != nil {
			return nil, fmt.Errorf("failed to scan answer: %w", err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (s *Service) SubmitTest(testID, userID int, birthDate string) (map[string]interface{}, error) {
	// 1. Проверяем существует ли тест
	var testTitle string
	err := s.db.QueryRow("SELECT title FROM tests WHERE id = ?", testID).Scan(&testTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	// 2. Определяем знак зодиака (id и название)
	zodiacID, zodiacName, err := ZodiacSign(birthDate)
	if err != nil {
		return nil, fmt.Errorf("failed to determine zodiac sign: %w", err)
	}

	// 3. Получаем результат по ID (он же id знака зодиака)
	result, err := s.getResultByID(testID, zodiacID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"zodiac_sign": zodiacName,
		"result": map[string]interface{}{
			"id":          result.ID,
			"title":       result.Title,
			"description": result.Description,
		},
	}, nil
}

func (s *Service) getResultByID(testID, resultID int) (*models.ZodiacResult, error) {
	var result models.ZodiacResult
	query := `SELECT id, title, description FROM results WHERE test_id = ? AND id = ?`
	err := s.db.QueryRow(query, testID, resultID).Scan(
		&result.ID, &result.Title, &result.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("result not found for id: %d", resultID)
		}
		return nil, fmt.Errorf("failed to get result: %w", err)
	}
	return &result, nil
}
func (s *Service) CalculateResult(testID int, answers map[int]int) (models.Result, error) {
	tagCount := make(map[string]int)

	// 1. Получаем result_tag для каждого ответа
	for _, answerID := range answers {
		var tag string

		err := s.db.QueryRow(
			"SELECT result_tag FROM answers WHERE id = ?",
			answerID,
		).Scan(&tag)

		if err != nil {
			continue // можно пропускать битые ответы
		}

		tagCount[tag]++
	}

	if len(tagCount) == 0 {
		return models.Result{}, errors.New("no valid answers")
	}

	// 2. Ищем самый частый тег
	var bestTag string
	max := 0

	for tag, count := range tagCount {
		if count > max {
			max = count
			bestTag = tag
		}
	}

	// 3. Получаем результат из БД
	var result models.Result

	err := s.db.QueryRow(`
		SELECT id, title, description
		FROM results
		WHERE test_id = ? AND result_tag = ?
	`, testID, bestTag).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
	)

	if err != nil {
		return models.Result{}, err
	}

	return result, nil
}

func (s *Service) CreateUser(name, password, birth, email string) (models.User, error) {
	var user models.User

	// пока без хеширования (для простоты)
	res, err := s.db.Exec(`
		INSERT INTO users (username, email, password, birth_date)
		VALUES (?, ?, ?, ?)
	`, name, birth, email, password)

	if err != nil {
		fmt.Print(err)
		return user, err
	}

	id, _ := res.LastInsertId()

	user = models.User{
		ID:    int(id),
		Name:  name,
		Birth: birth,
		Email: email,
	}

	return user, nil
}

// CheckFullMatch проверяет полное совпадение всех полей
func (s *Service) CheckFullMatch(name, password, birth, email string) (map[string]interface{}, error) {
	var user struct {
		ID       int
		Username string
		Email    string
		Birth    string
	}

	query := `SELECT id, username, email, COALESCE(birth_date, '') 
	          FROM users 
	          WHERE username = ? AND password = ? AND email = ? AND birth_date = ?`

	err := s.db.QueryRow(query, name, password, email, birth).Scan(&user.ID, &user.Username, &user.Email, &user.Birth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]interface{}{"found": false}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"found":     true,
		"all_match": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Username,
			"email": user.Email,
			"birth": user.Birth,
		},
	}, nil
}

// CheckEmailExists проверяет, существует ли пользователь с таким email
func (s *Service) CheckEmailExists(email string) (map[string]interface{}, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	err := s.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"found": exists}, nil
}

// CheckPasswordOwner находит пользователя (кроме указанного email), у которого такой пароль
func (s *Service) CheckPasswordOwner(password, excludeEmail string) (map[string]interface{}, error) {
	var owner struct {
		Username string
		Email    string
	}

	query := `SELECT username, email FROM users WHERE password = ? AND email != ? LIMIT 1`
	err := s.db.QueryRow(query, password, excludeEmail).Scan(&owner.Username, &owner.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]interface{}{"found": false}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"found":           true,
		"suggested_name":  owner.Username,
		"suggested_email": owner.Email,
	}, nil
}

func (s *Service) CheckEmailAndPassword(email, password string) (map[string]interface{}, error) {
	var user struct {
		ID    int
		Name  string
		Birth string
	}
	query := `SELECT id, username, COALESCE(birth_date, '') FROM users WHERE email = ? AND password = ?`
	err := s.db.QueryRow(query, email, password).Scan(&user.ID, &user.Name, &user.Birth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]interface{}{"found": false}, nil
		}
		return nil, err
	}
	return map[string]interface{}{
		"found": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"birth": user.Birth,
		},
	}, nil
}
