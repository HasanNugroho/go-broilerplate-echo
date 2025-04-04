package users

import (
	"errors"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/users/entity"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *config.DatabaseConfig
}

func NewUserRepository(db *config.DatabaseConfig) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx *gin.Context, user *entity.User) error {
	result := u.db.Client.Create(&user)
	return result.Error
}

func (u *UserRepository) FindByEmail(ctx *gin.Context, email string) (entity.User, error) {
	var user entity.User
	result := u.db.Client.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindById(ctx *gin.Context, id string) (entity.User, error) {
	var user entity.User
	result := u.db.Client.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, gorm.ErrRecordNotFound
		}
		return entity.User{}, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error) {
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
	var userModels []UserModelResponse
	for _, user := range users {
		userModels = append(userModels, UserModelResponse{
			ID:        (user.ID).String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		})
	}

	return userModels, int(totalItems), nil
}

func (u *UserRepository) Update(ctx *gin.Context, id string, user *entity.User) error {
	return u.db.Client.Where("id = ?", id).Updates(user).Error
}

func (u *UserRepository) Delete(ctx *gin.Context, id string) error {
	return u.db.Client.Where("id", id).Delete(&entity.User{}).Error
}
