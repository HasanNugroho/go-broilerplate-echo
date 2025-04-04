package users

import (
	"errors"
	"fmt"

	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
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
func (c *UserHandler) Create(ctx *gin.Context) {
	var user UserCreateModel
	ctx.Bind(&user)
	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return
	}

	if err := c.userService.Create(ctx, &user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, 201, "users created successfully", nil)
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
func (c *UserHandler) FindAll(ctx *gin.Context) {
	var filter shared.PaginationFilter

	// Binding query parameters
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.SendError(ctx, 400, "failed to fetch users", err)
		return
	}

	users, err := c.userService.FindAll(ctx, &filter)
	if err != nil {
		utils.SendError(ctx, 400, "failed to fetch users", err)
		return
	}
	utils.SendSuccess(ctx, 200, "users retrieved successfully", users)
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
func (c *UserHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	validate := validator.New()
	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	user, err := c.userService.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, 404, fmt.Sprintf("User with ID %s not found", id), err.Error())
			return
		}
		utils.SendError(ctx, 500, "Failed to fetch user", err.Error())
		return
	}
	utils.SendSuccess(ctx, 200, "User retrieved successfully", user)
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
func (c *UserHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var user UserUpdateModel
	validate := validator.New()

	ctx.Bind(&user)

	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, 400, "Bad request", err.Error())
		return
	}

	if err := c.userService.Update(ctx, id, &user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err)
		return
	}

	utils.SendSuccess(ctx, 201, "users updated successfully", nil)
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
func (c *UserHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	validate := validator.New()

	if err := validate.Var(id, "required,ulid"); err != nil {
		utils.SendError(ctx, 400, "Invalid ID", "ID is not a valid ULID")
		return
	}

	err := c.userService.Delete(ctx, id)
	if err != nil {
		utils.SendError(ctx, 400, "failed to delete user", err)
		return
	}
	utils.SendSuccess(ctx, 200, "user deleted successfully", nil)
}
