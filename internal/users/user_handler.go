package users

import (
	"github.com/HasanNugroho/starter-golang/internal/pkg/response"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
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
// @Param        user  body  model.UserCreateUpdateModel  true  "User Data"
// @Success      201  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users [post]
func (c *UserHandler) Create(ctx *gin.Context) {
	var user model.UserCreateUpdateModel
	ctx.Bind(&user)
	if err := c.userService.Create(ctx, &user); err != nil {
		response.SendError(ctx, 400, "create data failed", err)
	}

	response.SendSuccess(ctx, 201, "users created successfully", nil)
}

// FindAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   response.SuccessResponse{data=[]model.UserModelResponse}
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users [get]
func (c *UserHandler) FindAll(ctx *gin.Context) {
	users, err := c.userService.FindAll(ctx)
	if err != nil {
		response.SendError(ctx, 500, "failed to fetch users", err)
		return
	}
	response.SendSuccess(ctx, 200, "users retrieved successfully", users)
}
