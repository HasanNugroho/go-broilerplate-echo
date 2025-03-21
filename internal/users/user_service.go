package users

import (
	"github.com/HasanNugroho/starter-golang/internal/pkg/security"
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/HasanNugroho/starter-golang/internal/users/repository"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repo repository.IUserRepository
}

// ✅ Pastikan constructor mengembalikan *UserService
func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// ✅ Implementasi harus cocok dengan `IUserService`
func (u *UserService) Create(ctx *gin.Context, user *model.UserCreateUpdateModel) error {
	password, err := security.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}

	payload := entity.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: password,
	}
	return u.repo.Create(ctx, &payload)
}

func (u *UserService) FindById(ctx *gin.Context, id string) (model.UserModel, error) {
	// Implementasi yang sesuai
	panic("not implemented")
}

func (u *UserService) FindAll(ctx *gin.Context) ([]model.UserModel, error) {
	users, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert entity.User ke model.UserModel
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

func (u *UserService) Update(ctx *gin.Context, id string, user model.UserCreateUpdateModel) error {
	// Implementasi yang sesuai
	panic("not implemented")
}

func (u *UserService) Delete(ctx *gin.Context, id string) error {
	// Implementasi yang sesuai
	panic("not implemented")
}
