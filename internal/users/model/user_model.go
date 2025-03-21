package model

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
)

type UserModel struct {
	ID        utils.ULID `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserCreateUpdateModel struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
