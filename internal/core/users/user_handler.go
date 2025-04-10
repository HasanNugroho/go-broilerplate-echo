package users

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(us IUserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// CreateUser godoc
// @Summary      Create an user
// @Description  Create an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  UserCreateModel  true  "User Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /users [post]
// @Security ApiKeyAuth
func (c *UserHandler) Create(ctx echo.Context) error {
	var user UserCreateModel
	ctx.Bind(&user)
	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return err
	}

	if err := c.userService.Create(ctx, &user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return err
	}

	utils.SendSuccess(ctx, 201, "users created successfully", nil)
	return nil
}

// FindAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param limit query int false "total data per-page" minimum(1) default(10)
// @Param page query int false "page" minimum(1) default(1)
// @Param search query string false "keyword"
// @Success      200     {object}  shared.Response{data=shared.DataWithPagination{items=[]UserModelResponse}}
// @Failure      500     {object}  shared.Response
// @Router       /users [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindAll(ctx echo.Context) error {
	var filter shared.PaginationFilter

	// Binding query parameters
	if err := ctx.Bind(&filter); err != nil {
		utils.SendError(ctx, 400, "failed to fetch users", err)
		return err
	}

	users, err := c.userService.FindAll(ctx, &filter)
	if err != nil {
		utils.SendError(ctx, 400, "failed to fetch users", err)
		return err
	}

	utils.SendSuccess(ctx, 200, "users retrieved successfully", users)
	return nil
}

// FindUser godoc
// @Summary      Get all users
// @Description  Retrieve a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  shared.Response{data=UserModel}
// @Failure      500     {object}  shared.Response
// @Router       /users/{id} [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindById(ctx echo.Context) error {
	id := ctx.Param("id")

	validate := validator.New()
	if err := validate.Var(id, "required"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", nil)
		return err
	}

	user, err := c.userService.FindById(ctx, id)
	if err != nil {
		// utils.SendError(ctx, 404, fmt.Sprintf("User with ID %s not found", id), err.Error())
		// return err
		utils.SendError(ctx, 500, "Failed to fetch user", err.Error())
		return err
	}
	utils.SendSuccess(ctx, 200, "User retrieved successfully", user)
	return nil
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Param        user  body  UserUpdateModel  true  "User Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /users/{id} [put]
// @Security ApiKeyAuth
func (c *UserHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var user UserUpdateModel
	validate := validator.New()

	ctx.Bind(&user)

	if err := validate.Var(id, "required"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", nil)
		return err
	}

	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, 400, "Bad request", err.Error())
		return err
	}

	if err := c.userService.Update(ctx, id, &user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err)
		return err
	}

	utils.SendSuccess(ctx, 201, "users updated successfully", nil)
	return nil
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  shared.Response
// @Failure      500     {object}  shared.Response
// @Router       /users/{id} [delete]
// @Security ApiKeyAuth
func (c *UserHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	validate := validator.New()

	if err := validate.Var(id, "required"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", nil)
		return err
	}

	err := c.userService.Delete(ctx, id)
	if err != nil {
		utils.SendError(ctx, 400, "failed to delete user", err)
		return err
	}
	utils.SendSuccess(ctx, 200, "user deleted successfully", nil)
	return nil
}
