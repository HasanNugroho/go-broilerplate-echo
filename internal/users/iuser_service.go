package users

import (
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type IUserService interface {
	Create(ctx *gin.Context, user *model.UserCreateUpdateModel) error
	FindById(ctx *gin.Context, id string) (model.UserModel, error)
	FindAll(ctx *gin.Context) ([]model.UserModel, error)
	Update(ctx *gin.Context, id string, user model.UserCreateUpdateModel) error
	Delete(ctx *gin.Context, id string) error
}
