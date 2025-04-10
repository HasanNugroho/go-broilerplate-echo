package roles

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	roleService IRoleService
	validate    *validator.Validate
}

func NewRoleHandler(rs IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: rs,
		validate:    validator.New(),
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
func (c *RoleHandler) Create(ctx echo.Context) error {
	var role RoleUpdateModel
	ctx.Bind(&role)

	if err := c.validate.Struct(role); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	if err := c.roleService.Create(ctx, &role); err != nil {
		return err
	}

	utils.SendSuccess(ctx, 201, "roles created successfully", nil)
	return nil
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
func (c *RoleHandler) FindAll(ctx echo.Context) error {
	var filter shared.PaginationFilter

	if err := ctx.Bind(&filter); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	roles, err := c.roleService.FindAll(ctx, &filter)
	if err != nil {
		return err
	}

	utils.SendSuccess(ctx, 200, "roles retrieved successfully", roles)
	return nil
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
func (c *RoleHandler) FindById(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.validate.Var(id, "required"); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	role, err := c.roleService.FindById(ctx, id)
	if err != nil {
		return err
	}

	utils.SendSuccess(ctx, 200, "role retrieved successfully", role)
	return nil
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
func (c *RoleHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var role RoleUpdateModel

	ctx.Bind(&role)

	if err := c.validate.Var(id, "required"); err != nil {
		return utils.NewBadRequest("Invalid ID")
	}

	if err := c.validate.Struct(role); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	if err := c.roleService.Update(ctx, id, &role); err != nil {
		return err
	}

	utils.SendSuccess(ctx, 201, "roles updated successfully", nil)
	return nil
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
func (c *RoleHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.validate.Var(id, "required"); err != nil {
		return utils.NewBadRequest("Invalid ID")
	}

	err := c.roleService.Delete(ctx, id)
	if err != nil {
		return err
	}

	utils.SendSuccess(ctx, 200, "role deleted successfully", nil)
	return nil
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
func (c *RoleHandler) AssignUser(ctx echo.Context) error {
	var payload AssignRoleModel
	ctx.Bind(&payload)

	if err := c.validate.Struct(payload); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	if err := c.roleService.AssignUser(ctx, &payload); err != nil {
		return err
	}

	utils.SendSuccess(ctx, 201, "Assign user successfully", nil)
	return nil
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
// @Router       /roles/unassign [post]
// @Security ApiKeyAuth
func (c *RoleHandler) UnAssignUser(ctx echo.Context) error {
	var payload AssignRoleModel
	ctx.Bind(&payload)

	if err := c.validate.Struct(payload); err != nil {
		return utils.NewBadRequest(err.Error())
	}

	if err := c.roleService.UnassignUser(ctx, &payload); err != nil {
		return err
	}

	utils.SendSuccess(ctx, 201, "UnAssign user successfully", nil)
	return nil
}
