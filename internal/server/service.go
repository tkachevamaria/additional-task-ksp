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

// SubmitTest сохраняет результаты прохождения теста и возвращает результат
func (s *Service) SubmitTest(testID, userID int, answers map[int]int) (*models.TestResult, error) {
	// 1. Проверяем существует ли тест
	var testTitle string
	err := s.db.QueryRow("SELECT title FROM tests WHERE id = ?", testID).Scan(&testTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	// 2. Проверяем все ли ответы валидны
	resultTags, err := s.validateAnswers(answers)
	if err != nil {
		return nil, err
	}

	// 3. Подсчитываем результат (какой result_tag чаще всего встречается)
	finalTag := s.calculateResultTag(resultTags)

	// 4. Получаем информацию о результате из таблицы results
	result, err := s.getResultByTag(testID, finalTag)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// validateAnswers проверяет, что все answers существуют и возвращает их result_tag
func (s *Service) validateAnswers(answers map[int]int) ([]string, error) {
	var resultTags []string

	for questionID, answerID := range answers {
		var resultTag string
		query := "SELECT result_tag FROM answers WHERE id = ? AND question_id = ?"
		err := s.db.QueryRow(query, answerID, questionID).Scan(&resultTag)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("invalid answer: question_id=%d, answer_id=%d", questionID, answerID)
			}
			return nil, fmt.Errorf("failed to validate answer: %w", err)
		}
		resultTags = append(resultTags, resultTag)
	}

	return resultTags, nil
}

// calculateResultTag определяет финальный тег результата (тот, который чаще всего встречается)
func (s *Service) calculateResultTag(resultTags []string) string {
	if len(resultTags) == 0 {
		return "unknown"
	}

	tagCount := make(map[string]int)
	for _, tag := range resultTags {
		tagCount[tag]++
	}

	var maxTag string
	maxCount := 0
	for tag, count := range tagCount {
		if count > maxCount {
			maxCount = count
			maxTag = tag
		}
	}

	return maxTag
}

// getResultByTag получает результат из таблицы results по test_id и result_tag
func (s *Service) getResultByTag(testID int, resultTag string) (*models.TestResult, error) {
	var result models.TestResult
	query := `
        SELECT id, title, description, result_tag 
        FROM results 
        WHERE test_id = ? AND result_tag = ?
    `
	err := s.db.QueryRow(query, testID, resultTag).Scan(
		&result.ID, &result.Title, &result.Description, &result.ResultTag,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("result not found for tag: %s", resultTag)
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

	// ⚠️ пока без хеширования (для простоты)
	res, err := s.db.Exec(`
		INSERT INTO users (username, email, password_hash, birth_date)
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
