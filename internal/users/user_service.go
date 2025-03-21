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
	return u.repo.Create(&payload, ctx)
}

func (u *UserService) FindById(ctx *gin.Context, id string) (model.UserModel, error) {
	// Implementasi yang sesuai
	panic("not implemented")
}

func (u *UserService) FindAll(ctx *gin.Context, search interface{}) ([]model.UserModel, error) {
	// Implementasi yang sesuai
	panic("not implemented")
}

func (u *UserService) Update(ctx *gin.Context, id string, user model.UserCreateUpdateModel) error {
	// Implementasi yang sesuai
	panic("not implemented")
}

func (u *UserService) Delete(ctx *gin.Context, id string) error {
	// Implementasi yang sesuai
	panic("not implemented")
}
