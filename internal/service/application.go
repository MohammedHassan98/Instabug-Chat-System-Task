package service

import (
	"chat-system/internal/db/models"
	"context"
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"
)

type ApplicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, name string) (*models.Application, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	app := models.Application{
		Name:  name,
		Token: token,
	}

	if err := s.db.Create(&app).Error; err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *ApplicationService) GetApplicationByToken(ctx context.Context, token string) (*models.Application, error) {
	var app models.Application
	if err := s.db.Where("token = ?", token).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *ApplicationService) GetAllApplications(ctx context.Context) ([]models.Application, error) {
	var apps []models.Application

	if err := s.db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, token string, name string) (*models.Application, error) {
	app, err := s.GetApplicationByToken(ctx, token)

	if err != nil {
		return nil, err
	}

	app.Name = name
	if err := s.db.Save(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}
