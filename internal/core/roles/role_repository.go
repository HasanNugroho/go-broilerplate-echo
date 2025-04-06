package roles

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type RoleRepository struct {
	app *app.Apps
}

func NewRoleRepository(app *app.Apps) *RoleRepository {
	return &RoleRepository{
		app: app,
	}
}

func (app *RoleRepository) Create(ctx *gin.Context, role *Role) error {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) FindById(ctx *gin.Context, id string) (Role, error) {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error) {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) Update(ctx *gin.Context, id string, role *Role) error {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) Delete(ctx *gin.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) AssignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}

func (app *RoleRepository) UnassignUser(ctx *gin.Context, userId string, roleId string) error {
	panic("not implemented") // TODO: Implement
}
