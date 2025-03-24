package users

import (
	"errors"
	"fmt"

	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/HasanNugroho/starter-golang/internal/users/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Create(ctx *gin.Context, user *model.UserCreateModel) error {
	existingUser, err := u.repo.FindByEmail(ctx, user.Email)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if existingUser.Email != "" {
		return fmt.Errorf("user with Email %s already exists", user.Email)
	}

	password, err := utils.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}

	payload := entity.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: password,
	}

	err = u.repo.Create(ctx, &payload)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) FindById(ctx *gin.Context, id string) (model.UserModel, error) {
	user, err := u.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserModel{}, fmt.Errorf("user with ID %s not found", id)
		}
		return model.UserModel{}, err
	}

	return model.UserModel{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserService) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error) {
	users, totalItems, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return shared.DataWithPagination{}, err
	}

	// build pagination meta
	paginate := utils.BuildPagination(filter, int64(totalItems))

	// Buat response dengan pagination
	result := shared.DataWithPagination{
		Items:  users,
		Paging: paginate,
	}

	return result, nil
}

func (u *UserService) Update(ctx *gin.Context, id string, user *model.UserUpdateModel) error {
	existingUser, err := u.repo.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	updatedUser := entity.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: existingUser.Password,
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword([]byte(user.Password))
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		updatedUser.Password = hashedPassword
	}

	if err := u.repo.Update(ctx, id, &updatedUser); err != nil {
		return fmt.Errorf("failed to update user with ID %s: %w", id, err)
	}

	return nil
}

func (u *UserService) Delete(ctx *gin.Context, id string) error {
	if _, err := u.repo.FindById(ctx, id); err != nil {
		return fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user with ID %s: %w", id, err)
	}

	return nil
}
