package users

import (
	"github.com/HasanNugroho/starter-golang/internal/core/users/entity"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type IUserRepository interface {
	Create(ctx *gin.Context, user *entity.User) error
	FindByEmail(ctx *gin.Context, email string) (entity.User, error)
	FindById(ctx *gin.Context, id string) (entity.User, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error)
	Update(ctx *gin.Context, id string, user *entity.User) error
	Delete(ctx *gin.Context, id string) error
}

type IUserService interface {
	Create(ctx *gin.Context, user *UserCreateModel) error
	FindById(ctx *gin.Context, id string) (UserModel, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx *gin.Context, id string, user *UserUpdateModel) error
	Delete(ctx *gin.Context, id string) error
}
