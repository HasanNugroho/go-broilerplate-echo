package roles

import (
	"errors"
	"fmt"

	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RoleHandler struct {
	roleService IRoleService
}

func NewRoleHandler(rs IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: rs,
	}
}

// Createrole godoc
// @Summary      Create an role
// @Description  Create an role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        role  body  RoleUpdateModel  true  "role Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /roles [post]
// @Security ApiKeyAuth
func (c *RoleHandler) Create(ctx *gin.Context) {
	var role RoleUpdateModel
	ctx.Bind(&role)
	validate := validator.New()

	if err := validate.Struct(role); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return
	}

	if err := c.roleService.Create(ctx, &role); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, 201, "roles created successfully", nil)
}

// FindAllroles godoc
// @Summary      Get all roles
// @Description  Retrieve a list of all roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param limit query int false "total data per-page" minimum(1) default(10)
// @Param page query int false "page" minimum(1) default(1)
// @Param search query string false "keyword"
// @Success      200     {object}  shared.Response{data=shared.DataWithPagination{items=[]RoleModel}}
// @Failure      500     {object}  shared.Response
// @Router       /roles [get]
// @Security ApiKeyAuth
func (c *RoleHandler) FindAll(ctx *gin.Context) {
	var filter shared.PaginationFilter

	// Binding query parameters
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.SendError(ctx, 400, "failed to fetch roles", err)
		return
	}

	roles, err := c.roleService.FindAll(ctx, &filter)
	if err != nil {
		utils.SendError(ctx, 400, "failed to fetch roles", err)
		return
	}
	utils.SendSuccess(ctx, 200, "roles retrieved successfully", roles)
}

// Findrole godoc
// @Summary      Get all roles
// @Description  Retrieve a role by ID
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  shared.Response{data=RoleModel}
// @Failure      500     {object}  shared.Response
// @Router       /roles/{id} [get]
// @Security ApiKeyAuth
func (c *RoleHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	validate := validator.New()
	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	role, err := c.roleService.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, 404, fmt.Sprintf("role with ID %s not found", id), err.Error())
			return
		}
		utils.SendError(ctx, 500, "Failed to fetch role", err.Error())
		return
	}
	utils.SendSuccess(ctx, 200, "role retrieved successfully", role)
}

// Updaterole godoc
// @Summary      Update role
// @Description  Update role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Param        role  body  RoleUpdateModel  true  "role Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /roles/{id} [put]
// @Security ApiKeyAuth
func (c *RoleHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var role RoleUpdateModel
	validate := validator.New()

	ctx.Bind(&role)

	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	if err := validate.Struct(role); err != nil {
		utils.SendError(ctx, 400, "Bad request", err.Error())
		return
	}

	if err := c.roleService.Update(ctx, id, &role); err != nil {
		utils.SendError(ctx, 400, "create data failed", err)
		return
	}

	utils.SendSuccess(ctx, 201, "roles updated successfully", nil)
}

// Deleterole godoc
// @Summary      Delete role
// @Description  Delete role by ID
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  shared.Response
// @Failure      500     {object}  shared.Response
// @Router       /roles/{id} [delete]
// @Security ApiKeyAuth
func (c *RoleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	validate := validator.New()

	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	err := c.roleService.Delete(ctx, id)
	if err != nil {
		utils.SendError(ctx, 400, "failed to delete role", err)
		return
	}
	utils.SendSuccess(ctx, 200, "role deleted successfully", nil)
}

// Assignrole godoc
// @Summary      Assign an role
// @Description  Assign an role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        role  body  AssignRoleModel  true  "role Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /roles/assign [post]
// @Security ApiKeyAuth
func (c *RoleHandler) AssignUser(ctx *gin.Context) {
	var payload AssignRoleModel
	ctx.Bind(&payload)
	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		utils.SendError(ctx, 400, "Assign user failed", err.Error())
		return
	}

	if err := c.roleService.AssignUser(ctx, &payload); err != nil {
		utils.SendError(ctx, 400, "Assign user failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, 201, "Assign user successfully", nil)
}

// UnAssignrole godoc
// @Summary      UnAssign an role
// @Description  UnAssign an role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        role  body  AssignRoleModel  true  "role Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /roles/Unassign [post]
// @Security ApiKeyAuth
func (c *RoleHandler) UnAssignUser(ctx *gin.Context) {
	var payload AssignRoleModel
	ctx.Bind(&payload)
	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		utils.SendError(ctx, 400, "UnAssign user failed", err.Error())
		return
	}

	if err := c.roleService.UnassignUser(ctx, &payload); err != nil {
		utils.SendError(ctx, 400, "UnAssign user failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, 201, "UnAssign user successfully", nil)
}
