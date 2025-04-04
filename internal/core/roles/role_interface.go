package roles

import (
	"github.com/HasanNugroho/starter-golang/internal/core/roles/entity"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type IRoleRepository interface {
	Create(ctx gin.Context, role *entity.Role) error
	FindById(ctx *gin.Context, id string) (entity.Role, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error)
	Update(ctx *gin.Context, id string, role *entity.Role) error
	Delete(ctx *gin.Context, id string) error
}

type IRoleService interface {
	Create(ctx *gin.Context, user *RoleModel) error
	FindById(ctx *gin.Context, id string) (RoleModel, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx *gin.Context, id string, user *RoleUpdateModel) error
	Delete(ctx *gin.Context, id string) error
}
