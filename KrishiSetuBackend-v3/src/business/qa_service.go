package business

import (
	"errors"
	"strings"

	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type QAService struct {
	questionRepo QuestionRepository
	answerRepo   AnswerRepository
}

func NewQAService(questionRepo QuestionRepository, answerRepo AnswerRepository) *QAService {
	return &QAService{
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
	}
}

func (s *QAService) CreateQuestion(userID uint, req dto.CreateQuestionDTO) (*models.Question, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("Question cannot be empty")
	}

	q := models.Question{
		Content: content,
		UserID:  userID,
	}

	if err := s.questionRepo.Create(&q); err != nil {
		return nil, errors.New("Failed to create question")
	}

	// Fetch with preload User
	question, err := s.questionRepo.FindByIDWithUser(q.ID)
	if err != nil {
		return &q, nil
	}
	return question, nil
}

func (s *QAService) GetCommunityQuestions(userID uint, search string) ([]models.Question, error) {
	return s.questionRepo.GetCommunityQuestions(userID, search)
}

func (s *QAService) GetMyQuestions(userID uint, search string) ([]models.Question, error) {
	return s.questionRepo.GetMyQuestions(userID, search)
}

func (s *QAService) CreateAnswer(userID uint, req dto.CreateAnswerDTO) (*models.Answer, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" || req.QuestionID == 0 {
		return nil, errors.New("Invalid answer data")
	}

	a := models.Answer{
		Content:    content,
		UserID:     userID,
		QuestionID: req.QuestionID,
	}

	if err := s.answerRepo.Create(&a); err != nil {
		return nil, errors.New("Failed to create answer")
	}

	// Fetch with preload User
	answer, err := s.answerRepo.FindByID(a.ID)
	if err != nil {
		return &a, nil
	}
	return answer, nil
}

func (s *QAService) GetMyRepliedQuestions(userID uint, search string) ([]models.Question, error) {
	return s.questionRepo.GetMyRepliedQuestions(userID, search)
}

func (s *QAService) GetAnswersByQuestion(questionID uint) ([]models.Answer, error) {
	return s.answerRepo.GetByQuestionID(questionID)
}

func (s *QAService) GetMyAnswers(userID uint) ([]models.Answer, error) {
	return s.answerRepo.GetByUserID(userID)
}

func (s *QAService) VoteQuestion(userID uint, req dto.VoteQuestionDTO) (string, error) {
	if req.QuestionID == 0 || (req.Type != "up" && req.Type != "down") {
		return "", errors.New("Invalid vote data")
	}

	question, err := s.questionRepo.FindByID(req.QuestionID)
	if err != nil {
		return "", errors.New("Question not found")
	}

	existing, err := s.questionRepo.GetVote(userID, req.QuestionID)

	// No previous vote
	if err != nil || existing == nil {
		newVote := models.QuestionVote{
			UserID:     userID,
			QuestionID: req.QuestionID,
			Type:       req.Type,
		}
		if err := s.questionRepo.CreateVote(&newVote); err != nil {
			return "", err
		}

		if req.Type == "up" {
			question.Upvotes++
		} else {
			question.Downvotes++
		}
		s.questionRepo.UpdateVotesCount(question.ID, question.Upvotes, question.Downvotes)

		return "Voted", nil
	}

	// Same vote -> remove
	if existing.Type == req.Type {
		if err := s.questionRepo.DeleteVote(existing); err != nil {
			return "", err
		}

		if req.Type == "up" && question.Upvotes > 0 {
			question.Upvotes--
		}
		if req.Type == "down" && question.Downvotes > 0 {
			question.Downvotes--
		}
		s.questionRepo.UpdateVotesCount(question.ID, question.Upvotes, question.Downvotes)

		return "Vote removed", nil
	}

	// Switch vote
	if existing.Type == "up" {
		if question.Upvotes > 0 {
			question.Upvotes--
		}
		question.Downvotes++
	} else {
		if question.Downvotes > 0 {
			question.Downvotes--
		}
		question.Upvotes++
	}

	existing.Type = req.Type
	if err := s.questionRepo.UpdateVote(existing); err != nil {
		return "", err
	}
	s.questionRepo.UpdateVotesCount(question.ID, question.Upvotes, question.Downvotes)

	return "Vote switched", nil
}

func (s *QAService) VoteAnswer(userID uint, req dto.VoteAnswerDTO) (string, error) {
	if req.AnswerID == 0 || (req.Type != "up" && req.Type != "down") {
		return "", errors.New("Invalid vote data")
	}

	answer, err := s.answerRepo.FindByID(req.AnswerID)
	if err != nil {
		return "", errors.New("Answer not found")
	}

	existing, err := s.answerRepo.GetVote(userID, req.AnswerID)

	// No previous vote
	if err != nil || existing == nil {
		newVote := models.AnswerVote{
			UserID:   userID,
			AnswerID: req.AnswerID,
			Type:     req.Type,
		}
		if err := s.answerRepo.CreateVote(&newVote); err != nil {
			return "", err
		}

		if req.Type == "up" {
			answer.Upvotes++
		} else {
			answer.Downvotes++
		}
		s.answerRepo.UpdateVotesCount(answer.ID, answer.Upvotes, answer.Downvotes)

		return "Voted", nil
	}

	// Same vote -> remove
	if existing.Type == req.Type {
		if err := s.answerRepo.DeleteVote(existing); err != nil {
			return "", err
		}

		if req.Type == "up" && answer.Upvotes > 0 {
			answer.Upvotes--
		}
		if req.Type == "down" && answer.Downvotes > 0 {
			answer.Downvotes--
		}
		s.answerRepo.UpdateVotesCount(answer.ID, answer.Upvotes, answer.Downvotes)

		return "Vote removed", nil
	}

	// Switch vote
	if existing.Type == "up" {
		if answer.Upvotes > 0 {
			answer.Upvotes--
		}
		answer.Downvotes++
	} else {
		if answer.Downvotes > 0 {
			answer.Downvotes--
		}
		answer.Upvotes++
	}

	existing.Type = req.Type
	if err := s.answerRepo.UpdateVote(existing); err != nil {
		return "", err
	}
	s.answerRepo.UpdateVotesCount(answer.ID, answer.Upvotes, answer.Downvotes)

	return "Vote switched", nil
}
