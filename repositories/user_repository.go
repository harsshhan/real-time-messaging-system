package repositories

import (
	"backend/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user models.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, is_online, last_seen)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(context.Background(), query, user.Name, user.Email, user.PasswordHash, user.IsOnline, user.LastSeen)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash,COALESCE(profile_picture, ''), created_at, is_online, last_seen
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.IsOnline,
		&user.LastSeen,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.UserResponse, error) {
	query := `SELECT id, name, COALESCE(profile_picture, '') AS profile_picture, is_online, last_seen FROM users`

	rows, err := r.DB.Query(context.Background(), query)
	if(err!=nil){
		return nil,err
	}

	defer rows.Close()

	var users []models.UserResponse
	
	for rows.Next(){
		var user models.UserResponse
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.ProfilePicture,
			&user.IsOnline,
			&user.LastSeen,
		)
		if(err!=nil){
			return nil,err
		}

		users = append(users, user)
	}
	return users, nil
}