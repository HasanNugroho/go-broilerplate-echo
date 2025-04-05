package roles

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type RoleService struct {
	repo IRoleRepository
}

func NewRoleService(repo IRoleRepository) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (repo *RoleService) Create(ctx *gin.Context, user *RoleModel) error {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) FindById(ctx *gin.Context, id string) (RoleModel, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) Update(ctx *gin.Context, id string, user *RoleUpdateModel) error {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) Delete(ctx *gin.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) AssignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}

func (repo *RoleService) UnassignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}
