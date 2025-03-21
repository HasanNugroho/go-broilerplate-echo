package repository

import (
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type IUserRepository interface {
	Create(ctx *gin.Context, user *entity.User) error
	FindById(ctx *gin.Context, id string) (model.UserModel, error)
	FindAll(ctx *gin.Context) ([]model.UserModel, error)
	Update(ctx *gin.Context, id string, user entity.User) error
	Delete(ctx *gin.Context, id string) error
}
