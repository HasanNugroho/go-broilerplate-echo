package users

import (
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/labstack/echo/v4"
)

type IUserRepository interface {
	Create(ctx echo.Context, user *entities.User) error
	FindByEmail(ctx echo.Context, email string) (UserModel, error)
	FindById(ctx echo.Context, id string) (UserModel, error)
	FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error)
	Update(ctx echo.Context, id string, user *entities.User) error
	Delete(ctx echo.Context, id string) error
}

type IUserService interface {
	Create(ctx echo.Context, user *UserCreateModel) error
	FindById(ctx echo.Context, id string) (UserModel, error)
	FindAll(ctx echo.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx echo.Context, id string, user *UserUpdateModel) error
	Delete(ctx echo.Context, id string) error
}
