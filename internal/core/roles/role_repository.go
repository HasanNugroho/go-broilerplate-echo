package roles

import (
	"encoding/json"
	"errors"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleRepository struct {
	app *app.Apps
}

func NewRoleRepository(app *app.Apps) *RoleRepository {
	return &RoleRepository{
		app: app,
	}
}

func (r *RoleRepository) Create(ctx echo.Context, role *entities.Role) error {
	result := r.app.DB.WithContext(ctx.Request().Context()).Create(&role)
	return result.Error
}

func (r *RoleRepository) FindById(ctx echo.Context, id string) (RoleModel, error) {
	var role entities.Role
	result := r.app.DB.WithContext(ctx.Request().Context()).Where("id = ?", id).First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return RoleModel{}, nil
		}
		return RoleModel{}, result.Error
	}

	var permission []string
	err := json.Unmarshal([]byte(role.Permissions), &permission)
	if err != nil {
		return RoleModel{}, err
	}

	return RoleModel{
		ID:          role.ID.String(),
		Name:        role.Name,
		Permissions: permission,
	}, nil
}

func (r *RoleRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error) {
	var roles []entities.Role
	var totalItems int64

	query := r.app.DB.WithContext(ctx.Request().Context())

	// Hitung total data sebelum pagination
	if err := query.Model(&entities.Role{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Query data dengan pagination
	result := query.Scopes(utils.Paginate(filter)).
		Select([]string{"id", "name", "permissions"}).
		Find(&roles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Konversi ke response model
	var roleModels []RoleModel
	for _, role := range roles {
		var permission []string
		err := json.Unmarshal([]byte(role.Permissions), &permission)
		if err != nil {
			continue
		}

		roleModels = append(roleModels, RoleModel{
			ID:          (role.ID).String(),
			Name:        role.Name,
			Permissions: permission,
		})
	}

	return roleModels, int(totalItems), nil
}

func (r *RoleRepository) Update(ctx echo.Context, id string, role *entities.Role) error {
	return r.app.DB.WithContext(ctx.Request().Context()).Where("id = ?", id).Updates(role).Error
}

func (r *RoleRepository) Delete(ctx echo.Context, id string) error {
	return r.app.DB.WithContext(ctx.Request().Context()).Where("id", id).Delete(&entities.Role{}).Error
}

func (r *RoleRepository) AssignUser(ctx echo.Context, userId string, roleId string) error {
	return r.app.DB.WithContext(ctx.Request().Context()).Table("user_roles").Create(map[string]interface{}{
		"user_id": userId,
		"role_id": roleId,
	}).Error
}

func (r *RoleRepository) UnassignUser(ctx echo.Context, userId string, roleId string) error {
	return r.app.DB.WithContext(ctx.Request().Context()).
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userId, roleId).
		Delete(nil).Error
}
