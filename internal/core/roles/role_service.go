package roles

import (
	"fmt"
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
		return fmt.Errorf("permissions not found")
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
	role, err := r.repo.FindById(ctx, id)
	if err != nil {
		return RoleModel{}, err
	}
	return role, nil
}

func (r *RoleService) FindAll(ctx echo.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error) {
	users, totalItems, err := r.repo.FindAll(ctx, filter)
	fmt.Println(totalItems)
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
		return fmt.Errorf("role with ID %s not found: %w", id, err)
	}

	updatedRole := entities.Role{
		Name: currentRole.Name,
	}

	if role.Name != "" {
		updatedRole.Name = role.Name
	}

	if role.Permissions != nil {
		if len(utils.Intersection(role.Permissions, r.app.Config.ModulePermissions)) < 1 {
			return fmt.Errorf("permissions not found")
		}
		updatedRole.Permissions = role.Permissions
	} else {
		updatedRole.Permissions = currentRole.Permissions
	}

	return r.repo.Update(ctx, id, &updatedRole)
}

func (r *RoleService) Delete(ctx echo.Context, id string) error {
	if _, err := r.repo.FindById(ctx, id); err != nil {
		return fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	if err := r.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user with ID %s: %w", id, err)
	}

	return nil
}

func (r *RoleService) AssignUser(ctx echo.Context, payload *AssignRoleModel) error {
	err := r.repo.AssignUser(ctx, payload.UserID, payload.RoleID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleService) UnassignUser(ctx echo.Context, payload *AssignRoleModel) error {
	err := r.repo.UnassignUser(ctx, payload.UserID, payload.RoleID)
	if err != nil {
		return err
	}
	return nil
}
