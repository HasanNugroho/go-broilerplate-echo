package users

import (
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/gin-gonic/gin"
)

type IUserRepository interface {
	Create(ctx *gin.Context, user *entities.User) error
	FindByEmail(ctx *gin.Context, email string) (entities.User, error)
	FindById(ctx *gin.Context, id string) (entities.User, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error)
	Update(ctx *gin.Context, id string, user *entities.User) error
	Delete(ctx *gin.Context, id string) error
}

type IUserService interface {
	Create(ctx *gin.Context, user *UserCreateModel) error
	FindById(ctx *gin.Context, id string) (UserModel, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx *gin.Context, id string, user *UserUpdateModel) error
	Delete(ctx *gin.Context, id string) error
}
