package repository

import (
	"github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type UserRepository struct {
	db *configs.RDBMSConfig
}

func NewUserRepository(db *configs.RDBMSConfig) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(user *entity.User, ctx *gin.Context) error {
	result := u.db.Client.Create(&user)
	return result.Error
}

func (u *UserRepository) FindById(id string, ctx *gin.Context) (model.UserModel, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) FindAll(search interface{}, ctx *gin.Context) ([]model.UserModel, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) Update(id string, user entity.User, ctx *gin.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) Delete(id string, ctx *gin.Context) error {
	panic("not implemented") // TODO: Implement
}
