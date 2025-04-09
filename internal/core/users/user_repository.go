package users

import (
	"errors"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserRepository struct {
	app *app.Apps
}

func NewUserRepository(app *app.Apps) *UserRepository {
	return &UserRepository{
		app: app,
	}
}

func (u *UserRepository) Create(ctx echo.Context, user *entities.User) error {
	result := u.app.DB.Create(&user)
	return result.Error
}

func (u *UserRepository) FindByEmail(ctx echo.Context, email string) (entities.User, error) {
	var user entities.User
	err := u.app.DB.WithContext(ctx.Request().Context()).
		Preload("Roles").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, nil
		}
		return entities.User{}, err
	}
	return user, nil
}

func (u *UserRepository) FindById(ctx echo.Context, id string) (entities.User, error) {
	var user entities.User
	result := u.app.DB.WithContext(ctx.Request().Context()).
		Preload("Roles").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.User{}, gorm.ErrRecordNotFound
		}
		return entities.User{}, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error) {
	var users []entities.User
	var totalItems int64

	query := u.app.DB.WithContext(ctx.Request().Context())

	// Hitung total data sebelum pagination
	if err := query.Model(&entities.User{}).Count(&totalItems).Error; err != nil {
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

func (u *UserRepository) Update(ctx echo.Context, id string, user *entities.User) error {
	return u.app.DB.WithContext(ctx.Request().Context()).Where("id = ?", id).Updates(user).Error
}

func (u *UserRepository) Delete(ctx echo.Context, id string) error {
	return u.app.DB.WithContext(ctx.Request().Context()).Where("id", id).Delete(&entities.User{}).Error
}
