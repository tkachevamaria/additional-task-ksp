package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GET /tests
func (h *Handler) GetAllTests(c *gin.Context) {
	tests, err := h.service.GetAllTests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tests)
}

// GET /tests/:id
func (h *Handler) GetTestByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	test, err := h.service.GetTestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

func (h *Handler) SubmitTest(c *gin.Context) {
	log.Printf("📨 [SubmitTest Handler] Получен запрос на отправку теста")
	log.Printf("📋 [SubmitTest Handler] Параметры запроса: %+v", c.Params)

	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("❌ [SubmitTest Handler] Невалидный ID теста: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test ID"})
		return
	}
	log.Printf("🆔 [SubmitTest Handler] Test ID: %d", testID)

	var req struct {
		UserID    int    `json:"user_id"`
		BirthDate string `json:"birth_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [SubmitTest Handler] Ошибка парсинга JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	log.Printf("📦 [SubmitTest Handler] Данные запроса: UserID=%d, BirthDate=%s", req.UserID, req.BirthDate)

	log.Printf("🔀 [SubmitTest Handler] Вызываю сервис SubmitTest")
	result, err := h.service.SubmitTest(testID, req.UserID, req.BirthDate)
	if err != nil {
		log.Printf("❌ [SubmitTest Handler] Ошибка сервиса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ [SubmitTest Handler] Успешный ответ: %+v", result)
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Birth    string `json:"birth"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if req.Name == "" {
		c.JSON(400, gin.H{"error": "name is required"})
		return
	}

	user, err := h.service.CreateUser(req.Name, req.Birth, req.Email, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) CheckPasswordOwner(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.service.CheckPasswordOwner(req.Password, req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *Handler) CheckFullMatch(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
		Birth    string `json:"birth"`
		Email    string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.service.CheckFullMatch(req.Name, req.Password, req.Birth, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CheckEmailExists(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.service.CheckEmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CheckEmailAndPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.service.CheckEmailAndPassword(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
