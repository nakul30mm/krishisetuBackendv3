package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *questionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) FindByID(id uint) (*models.Question, error) {
	var question models.Question
	if err := r.db.First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByIDWithUser(id uint) (*models.Question, error) {
	var question models.Question
	if err := r.db.Preload("User").First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) Update(question *models.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) GetCommunityQuestions(userID uint, search string) ([]models.Question, error) {
	var questions []models.Question

	query := r.db.Model(&models.Question{}).
		Select(`
		questions.id,
		questions.content,
		questions.user_id,
		questions.upvotes,
		questions.downvotes,
		questions.created_at,
		questions.updated_at,
		(SELECT COUNT(*) FROM answers 
		WHERE answers.question_id = questions.id) AS replies_count
		`).
		Where("questions.user_id != ?", userID)

	if search != "" {
		query = query.Where("questions.content LIKE ?", "%"+search+"%")
	}

	err := query.Preload("User").
		Order("questions.created_at DESC").
		Find(&questions).Error

	return questions, err
}

func (r *questionRepository) GetMyQuestions(userID uint, search string) ([]models.Question, error) {
	var questions []models.Question

	query := r.db.Model(&models.Question{}).
		Select(`
		questions.id,
		questions.content,
		questions.user_id,
		questions.upvotes,
		questions.downvotes,
		questions.created_at,
		questions.updated_at,
		(SELECT COUNT(*) FROM answers 
		WHERE answers.question_id = questions.id) AS replies_count
		`).
		Where("questions.user_id = ?", userID)

	if search != "" {
		query = query.Where("questions.content LIKE ?", "%"+search+"%")
	}

	err := query.Preload("User").
		Order("questions.created_at DESC").
		Find(&questions).Error

	return questions, err
}

func (r *questionRepository) GetMyRepliedQuestions(userID uint, search string) ([]models.Question, error) {
	var questions []models.Question

	query := r.db.Model(&models.Question{}).
		Select(`
			questions.id,
			questions.content,
			questions.user_id,
			questions.upvotes,
			questions.downvotes,
			questions.created_at,
			questions.updated_at,
			(SELECT COUNT(*) 
			 FROM answers 
			 WHERE answers.question_id = questions.id) AS replies_count
		`).
		Joins("JOIN answers ON answers.question_id = questions.id").
		Where("answers.user_id = ?", userID).
		Preload("User")

	if search != "" {
		query = query.Where("questions.content LIKE ?", "%"+search+"%")
	}

	err := query.Group("questions.id").
		Order("MAX(answers.created_at) DESC").
		Find(&questions).Error

	return questions, err
}

func (r *questionRepository) GetVote(userID uint, questionID uint) (*models.QuestionVote, error) {
	var vote models.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *questionRepository) CreateVote(vote *models.QuestionVote) error {
	return r.db.Create(vote).Error
}

func (r *questionRepository) DeleteVote(vote *models.QuestionVote) error {
	return r.db.Delete(vote).Error
}

func (r *questionRepository) UpdateVote(vote *models.QuestionVote) error {
	return r.db.Save(vote).Error
}

func (r *questionRepository) UpdateVotesCount(questionID uint, upvotes, downvotes int64) error {
	return r.db.Model(&models.Question{}).Where("id = ?", questionID).Updates(map[string]interface{}{
		"upvotes":   upvotes,
		"downvotes": downvotes,
	}).Error
}
