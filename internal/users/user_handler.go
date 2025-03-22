package users

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
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
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /users [post]
func (c *UserHandler) Create(ctx *gin.Context) {
	var user model.UserCreateUpdateModel
	ctx.Bind(&user)
	if err := c.userService.Create(ctx, &user); err != nil {
		utils.SendError(ctx, 400, "create data failed", err)
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
// @Success      200     {object}  shared.Response{data=shared.DataWithPagination{items=[]model.UserModelResponse}}
// @Failure      500     {object}  shared.Response
// @Router       /users [get]
func (c *UserHandler) FindAll(ctx *gin.Context) {
	var filter shared.PaginationFilter

	// Binding query parameters
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.SendError(ctx, 403, "failed to fetch users", err)
		return
	}

	users, err := c.userService.FindAll(ctx, &filter)
	if err != nil {
		utils.SendError(ctx, 403, "failed to fetch users", err)
		return
	}
	utils.SendSuccess(ctx, 200, "users retrieved successfully", users)
}
