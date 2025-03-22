package repository

import (
	"github.com/HasanNugroho/starter-golang/config"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type UserRepository struct {
	db *config.DBConfig
}

func NewUserRepository(db *config.DBConfig) *UserRepository {
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

func (u *UserRepository) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]model.UserModelResponse, int, error) {
	var users []entity.User
	var totalItems int64

	query := u.db.Client.WithContext(ctx)

	// Hitung total data sebelum pagination
	if err := query.Model(&entity.User{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Query data dengan pagination
	result := query.Scopes(utils.Paginate(filter)).
		Select([]string{"id", "name", "email", "created_at"}).
		Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Konversi ke response model
	var userModels []model.UserModelResponse
	for _, user := range users {
		userModels = append(userModels, model.UserModelResponse{
			ID:        (user.ID).String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		})
	}

	return userModels, int(totalItems), nil
}

func (u *UserRepository) Update(ctx *gin.Context, id string, user entity.User) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserRepository) Delete(ctx *gin.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
