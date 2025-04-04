package auth

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/auth/model"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	repo users.IUserRepository
}

func NewAuthService(repo users.IUserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Login(ctx *gin.Context, config *config.Config, email string, password string) (model.AuthResponse, error) {
	existingUser, err := a.repo.FindByEmail(ctx, email)

	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	if existingUser.Email == "" {
		return model.AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	if !utils.VerifyPassword(existingUser.Password, []byte(password)) {
		return model.AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	token, err := utils.GenerateAuthToken(config, users.UserModelResponse{
		ID:        (existingUser.ID).String(),
		Email:     existingUser.Email,
		Name:      existingUser.Name,
		CreatedAt: existingUser.CreatedAt,
	})

	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("Error creating token: %s", err.Error())
	}

	return token, nil
}

func (a *AuthService) Register(ctx *gin.Context, email string, password string) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AuthService) ResetPassword(ctx *gin.Context, email string, password string) error {
	panic("not implemented") // TODO: Implement
}
