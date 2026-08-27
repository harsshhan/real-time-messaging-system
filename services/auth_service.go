package services

import (
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"errors"
	"time"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) error {
	existingUser, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsOnline:     false,
		LastSeen:     time.Now(),
	}

	return s.UserRepo.Create(user)
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}	
	response := &models.LoginResponse{

		User: models.UserResponse{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			IsOnline:       user.IsOnline,
			LastSeen:       user.LastSeen,
		},
		Token: token,
	}

	return response, nil
}

