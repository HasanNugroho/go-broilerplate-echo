package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"unique"`
	Email      string    `json:"email" gorm:"unique"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastLogin  time.Time `json:"last_login"`
	IsActive   bool      `json:"is_active" gorm:"default:false"`
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Name      string    `json:"name"`
	LastLogin time.Time `json:"last_login"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id" gorm:"unique"`
	Email string    `json:"email" gorm:"unique"`
}
