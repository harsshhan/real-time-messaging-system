package services

import (
	"backend/models"
	"backend/repositories"
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	Redis *redis.Client
}

func NewUserService(userRepo *repositories.UserRepository, redis *redis.Client) *UserService {
	return &UserService{
		UserRepo: userRepo,
		Redis: redis,
	}
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {

	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for i := range users {
		key := "user:" + users[i].ID.String()
		exists, err := s.Redis.Exists(ctx, key).Result()
		if err != nil {

			return nil, err
		}
		users[i].IsOnline = exists > 0
	}

	return users, nil
}

func (s *UserService) UpdateLastSeen(userID uuid.UUID) error {
	return s.UserRepo.UpdateLastSeen(userID)
}