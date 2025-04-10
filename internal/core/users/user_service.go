package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Create(ctx echo.Context, user *UserCreateModel) error {
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

	payload := entities.User{
		Email:    user.Email,
		Name:     user.Name,
		Roles:    []bson.ObjectID{},
		Password: password,
	}

	if err = u.repo.Create(ctx, &payload); err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindById(ctx echo.Context, id string) (UserModel, error) {
	user, err := u.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrFileNotFound) {
			return UserModel{}, fmt.Errorf("user with ID %s not found", id)
		}
		return UserModel{}, err
	}

	return user, nil
}

func (u *UserService) FindAll(ctx echo.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error) {
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

func (u *UserService) Update(ctx echo.Context, id string, user *UserUpdateModel) error {
	existingUser, err := u.repo.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	updatedUser := entities.User{
		Email:     user.Email,
		Name:      user.Name,
		Password:  existingUser.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (u *UserService) Delete(ctx echo.Context, id string) error {
	if _, err := u.repo.FindById(ctx, id); err != nil {
		return fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user with ID %s: %w", id, err)
	}

	return nil
}
