package users

import (
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

func (c *UserHandler) Create(ctx *gin.Context) {
	var user model.UserCreateUpdateModel
	ctx.Bind(&user)
	if err := c.userService.Create(ctx, &user); err != nil {
		ctx.JSON(400, gin.H{
			"err": err,
		})
	}
}
