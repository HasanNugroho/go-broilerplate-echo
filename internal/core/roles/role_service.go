package roles

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
)

type RoleService struct {
	repo IRoleRepository
	app  *app.Apps
}

func NewRoleService(app *app.Apps, repo IRoleRepository) *RoleService {
	return &RoleService{
		repo: repo,
		app:  app,
	}
}

func (r *RoleService) Create(ctx echo.Context, user *RoleUpdateModel) error {
	if len(utils.Intersection(user.Permissions, r.app.Config.ModulePermissions)) < 1 {
		return utils.NewBadRequest("permissions not found")
	}

	payload := entities.Role{
		Name:        user.Name,
		Permissions: user.Permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.repo.Create(ctx, &payload); err != nil {
		return err
	}

	return nil
}

func (r *RoleService) FindById(ctx echo.Context, id string) (RoleModel, error) {
	return r.repo.FindById(ctx, id)
}

func (r *RoleService) FindAll(ctx echo.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error) {
	users, totalItems, err := r.repo.FindAll(ctx, filter)
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

func (r *RoleService) Update(ctx echo.Context, id string, role *RoleUpdateModel) error {
	currentRole, err := r.repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	updatedRole := entities.Role{
		Name: currentRole.Name,
	}

	if role.Name != "" {
		updatedRole.Name = role.Name
	}

	if role.Permissions != nil {
		if len(utils.Intersection(role.Permissions, r.app.Config.ModulePermissions)) < 1 {
			return utils.NewBadRequest("permissions not found")
		}
		updatedRole.Permissions = role.Permissions
	} else {
		updatedRole.Permissions = currentRole.Permissions
	}

	return r.repo.Update(ctx, id, &updatedRole)
}

func (r *RoleService) Delete(ctx echo.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

func (r *RoleService) AssignUser(ctx echo.Context, payload *AssignRoleModel) error {
	return r.repo.AssignUser(ctx, payload.UserID, payload.RoleID)
}

func (r *RoleService) UnassignUser(ctx echo.Context, payload *AssignRoleModel) error {
	return r.repo.UnassignUser(ctx, payload.UserID, payload.RoleID)
}
