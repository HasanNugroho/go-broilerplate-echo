package roles

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type IRoleRepository interface {
	Create(ctx *gin.Context, role *Role) error
	FindById(ctx *gin.Context, id string) (Role, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error)
	Update(ctx *gin.Context, id string, role *Role) error
	Delete(ctx *gin.Context, id string) error
	AssignUser(ctx *gin.Context, userId string, roleId string) error
	UnassignUser(ctx *gin.Context, userId string, roleId string) error
}

type IRoleService interface {
	Create(ctx *gin.Context, user *RoleModel) error
	FindById(ctx *gin.Context, id string) (RoleModel, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx *gin.Context, id string, user *RoleUpdateModel) error
	Delete(ctx *gin.Context, id string) error
	AssignUser(ctx *gin.Context, userId string, roleId string) error
	UnassignUser(ctx *gin.Context, userId string, roleId string) error
}
