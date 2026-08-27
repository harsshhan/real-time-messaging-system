package models

import ( "github.com/google/uuid" 
"time" )


type User struct {
    ID           uuid.UUID
    Name         string
    Email        string
    PasswordHash string
    ProfilePicture string
    CreatedAt    time.Time
    IsOnline     bool
    LastSeen     time.Time
}

type UserResponse struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
    Email string    `json:"email"`
    ProfilePicture string `json:"profile_picture"`
    IsOnline bool `json:"is_online"`
    LastSeen time.Time `json:"last_seen"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}