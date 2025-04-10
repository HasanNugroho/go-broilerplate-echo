package users

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	_, err := u.repo.FindByEmail(ctx, user.Email)

	if err != nil {
		return err
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
	return u.repo.FindById(ctx, id)
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
		return err
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
			return err
		}
		updatedUser.Password = hashedPassword
	}

	if err := u.repo.Update(ctx, id, &updatedUser); err != nil {
		return err
	}

	return nil
}

func (u *UserService) Delete(ctx echo.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
