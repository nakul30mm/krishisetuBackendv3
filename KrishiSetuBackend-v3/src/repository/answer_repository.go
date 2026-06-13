package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *answerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(answer *models.Answer) error {
	return r.db.Create(answer).Error
}

func (r *answerRepository) FindByID(id uint) (*models.Answer, error) {
	var answer models.Answer
	if err := r.db.First(&answer, id).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *answerRepository) Update(answer *models.Answer) error {
	return r.db.Save(answer).Error
}

func (r *answerRepository) GetByQuestionID(questionID uint) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.db.Preload("User").
		Where("question_id = ?", questionID).
		Order("created_at asc").
		Find(&answers).Error
	return answers, err
}

func (r *answerRepository) GetByUserID(userID uint) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&answers).Error
	return answers, err
}

func (r *answerRepository) GetVote(userID uint, answerID uint) (*models.AnswerVote, error) {
	var vote models.AnswerVote
	err := r.db.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *answerRepository) CreateVote(vote *models.AnswerVote) error {
	return r.db.Create(vote).Error
}

func (r *answerRepository) DeleteVote(vote *models.AnswerVote) error {
	return r.db.Delete(vote).Error
}

func (r *answerRepository) UpdateVote(vote *models.AnswerVote) error {
	return r.db.Save(vote).Error
}

func (r *answerRepository) UpdateVotesCount(answerID uint, upvotes, downvotes int64) error {
	return r.db.Model(&models.Answer{}).Where("id = ?", answerID).Updates(map[string]interface{}{
		"upvotes":   upvotes,
		"downvotes": downvotes,
	}).Error
}
