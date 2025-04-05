package users

import (
	"errors"

	"github.com/HasanNugroho/starter-golang/config"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *config.DatabaseConfig
}

func NewUserRepository(config *config.Config) *UserRepository {
	return &UserRepository{
		db: &config.DB,
	}
}

func (u *UserRepository) Create(ctx *gin.Context, user *User) error {
	result := u.db.Client.Create(&user)
	return result.Error
}

func (u *UserRepository) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var user User
	result := u.db.Client.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, nil
		}
		return User{}, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindById(ctx *gin.Context, id string) (User, error) {
	var user User
	result := u.db.Client.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, gorm.ErrRecordNotFound
		}
		return User{}, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error) {
	var users []User
	var totalItems int64

	query := u.db.Client.WithContext(ctx)

	// Hitung total data sebelum pagination
	if err := query.Model(&User{}).Count(&totalItems).Error; err != nil {
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

func (u *UserRepository) Update(ctx *gin.Context, id string, user *User) error {
	return u.db.Client.Where("id = ?", id).Updates(user).Error
}

func (u *UserRepository) Delete(ctx *gin.Context, id string) error {
	return u.db.Client.Where("id", id).Delete(&User{}).Error
}
