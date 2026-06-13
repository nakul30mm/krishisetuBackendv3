package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type QAHandler struct {
	Service *business.QAService
}

func NewQAHandler(service *business.QAService) *QAHandler {
	return &QAHandler{Service: service}
}

func (h *QAHandler) CreateQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateQuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	q, err := h.Service.CreateQuestion(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, q)
}

func (h *QAHandler) GetCommunityQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")
	search := c.Query("search")

	questions, err := h.Service.GetCommunityQuestions(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QAHandler) GetMyQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")
	search := c.Query("search")

	questions, err := h.Service.GetMyQuestions(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QAHandler) CreateAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateAnswerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	a, err := h.Service.CreateAnswer(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *QAHandler) GetMyRepliedQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")
	search := c.Query("search")

	questions, err := h.Service.GetMyRepliedQuestions(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QAHandler) GetAnswersByQuestion(c *gin.Context) {
	qidStr := c.Param("id")
	qid, err := strconv.ParseUint(qidStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	answers, err := h.Service.GetAnswersByQuestion(uint(qid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answers)
}

func (h *QAHandler) GetMyAnswers(c *gin.Context) {
	userID := c.GetUint("user_id")

	answers, err := h.Service.GetMyAnswers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answers)
}

func (h *QAHandler) VoteQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.VoteQuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	msg, err := h.Service.VoteQuestion(userID, req)
	if err != nil {
		if err.Error() == "Question not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *QAHandler) VoteAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.VoteAnswerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	msg, err := h.Service.VoteAnswer(userID, req)
	if err != nil {
		if err.Error() == "Answer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}
