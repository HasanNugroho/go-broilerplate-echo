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

func (u *UserRepository) Create(ctx *gin.Context, user *entity.User) error {
	result := u.db.Client.Create(&user)
	return result.Error
}

func (u *UserRepository) FindById(ctx *gin.Context, id string) (model.UserModel, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) FindAll(ctx *gin.Context) ([]model.UserModel, error) {
	var users []entity.User
	query := u.db.Client.WithContext(ctx)

	result := query.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	var userModels []model.UserModel
	for _, user := range users {
		userModels = append(userModels, model.UserModel{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	return userModels, nil
}

func (u *UserRepository) Update(ctx *gin.Context, id string, user entity.User) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) Delete(ctx *gin.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
