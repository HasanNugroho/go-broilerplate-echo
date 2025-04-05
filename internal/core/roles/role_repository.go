package roles

import (
	"github.com/HasanNugroho/starter-golang/config"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type RoleRepository struct {
	db *config.DatabaseConfig
}

func NewRoleRepository(config *config.Config) *RoleRepository {
	return &RoleRepository{
		db: &config.DB,
	}
}

func (db *RoleRepository) Create(ctx *gin.Context, role *Role) error {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) FindById(ctx *gin.Context, id string) (Role, error) {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error) {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) Update(ctx *gin.Context, id string, role *Role) error {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) Delete(ctx *gin.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) AssignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}

func (db *RoleRepository) UnassignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}
