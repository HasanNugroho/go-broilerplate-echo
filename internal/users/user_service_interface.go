package users

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/users/model"
	"github.com/gin-gonic/gin"
)

type IUserService interface {
	Create(ctx *gin.Context, user *model.UserCreateModel) error
	FindById(ctx *gin.Context, id string) (model.UserModel, error)
	FindAll(ctx *gin.Context, filter *shared.PaginationFilter) (shared.DataWithPagination, error)
	Update(ctx *gin.Context, id string, user *model.UserUpdateModel) error
	Delete(ctx *gin.Context, id string) error
}
