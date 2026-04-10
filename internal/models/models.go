package models

import "time"

type ErrorResponse struct {
	Message string `json:"message"` // tag json
}

// type User struct {
// 	ID           int
// 	Username     string
// 	Email        string
// 	PasswordHash string
// }

type Test struct {
	ID    int
	Title string
}

type TestDetail struct {
	ID        int              `json:"id"`
	Title     string           `json:"title"`
	Questions []QuestionDetail `json:"questions"`
}

type QuestionDetail struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	ResultTag string `json:"result_tag"`
}

type TestResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ResultTag   string `json:"result_tag"`
}

type UserResult struct {
	ID          int       `json:"id"`
	TestID      int       `json:"test_id"`
	TestTitle   string    `json:"test_title"`
	ResultTag   string    `json:"result_tag"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completed_at"`
}

type Result struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Birth string `json:"birth"`
	Email string `json:"email"`
}
