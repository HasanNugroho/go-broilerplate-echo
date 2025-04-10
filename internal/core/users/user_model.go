package users

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/core/roles"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserModel struct {
	ID        bson.ObjectID     `json:"id" bson:"_id"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Password  string            `json:"password,omitempty"`
	RolesData []roles.RoleModel `json:"roles_data" bson:"roles_data"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UserCreateModel struct {
	Email    string `json:"email" validate:"email"`
	Name     string `json:"name" validate:""`
	Password string `json:"password" validate:"min=6"`
}

type UserUpdateModel struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserModelResponse struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Email     string        `bson:"email" json:"email"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
