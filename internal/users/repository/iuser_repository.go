package repository

import (
	"github.com/HasanNugroho/starter-golang/internal/users/entity"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type IUserRepository interface {
	Create(user *entity.User, ctx *gin.Context) error
	FindById(id string, ctx *gin.Context) (model.UserModel, error)
	FindAll(search interface{}, ctx *gin.Context) ([]model.UserModel, error)
	Update(id string, user entity.User, ctx *gin.Context) error
	Delete(id string, ctx *gin.Context) error
}
