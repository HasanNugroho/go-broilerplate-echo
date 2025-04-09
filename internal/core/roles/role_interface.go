package roles

import (
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/labstack/echo/v4"
)

type IRoleRepository interface {
	Create(ctx echo.Context, role *entities.Role) error
	FindById(ctx echo.Context, id string) (RoleModel, error)
	FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error)
	Update(ctx echo.Context, id string, role *entities.Role) error
	Delete(ctx echo.Context, id string) error
	AssignUser(ctx echo.Context, userId string, roleId string) error
	UnassignUser(ctx echo.Context, userId string, roleId string) error
}

type IRoleService interface {
	Create(ctx echo.Context, user *RoleUpdateModel) error
	FindById(ctx echo.Context, id string) (RoleModel, error)
	FindAll(ctx echo.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx echo.Context, id string, user *RoleUpdateModel) error
	Delete(ctx echo.Context, id string) error
	AssignUser(ctx echo.Context, payload *AssignRoleModel) error
	UnassignUser(ctx echo.Context, payload *AssignRoleModel) error
}
