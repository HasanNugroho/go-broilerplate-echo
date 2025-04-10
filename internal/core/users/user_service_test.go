package users_test

import (
	"errors"
	"testing"

	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByEmail(ctx echo.Context, email string) (users.UserModel, error) {
	panic("not implemented")
}

func (m *MockUserRepo) Update(ctx echo.Context, id string, user *entities.User) error {
	panic("not implemented")
}

func (m *MockUserRepo) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]users.UserModelResponse, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]users.UserModelResponse), args.Int(1), args.Error(2)
}

func (m *MockUserRepo) FindById(ctx echo.Context, id string) (users.UserModel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(users.UserModel), args.Error(1)
}

func (m *MockUserRepo) Delete(ctx echo.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) Create(ctx echo.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(1)
}

func newTestContext() echo.Context {
	return echo.New().NewContext(nil, nil)
}

func newPaginationFilter() *shared.PaginationFilter {
	return &shared.PaginationFilter{Limit: 10, Page: 1}
}

func TestUserService_FindAll_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := users.NewUserService(mockRepo)
	ctx := newTestContext()
	filter := newPaginationFilter()

	expectedUsers := []users.UserModelResponse{
		{ID: bson.NewObjectID(), Name: "John Doe", Email: "john@example.com"},
	}
	expectedTotal := 1

	mockRepo.
		On("FindAll", ctx, filter).
		Return(expectedUsers, expectedTotal, nil)

	result, err := service.FindAll(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(expectedTotal), result.Paging.TotalItems)
	assert.Equal(t, expectedUsers, result.Items)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindAll_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := users.NewUserService(mockRepo)
	ctx := newTestContext()
	filter := newPaginationFilter()

	mockRepo.
		On("FindAll", ctx, filter).
		Return([]users.UserModelResponse{}, 0, errors.New("db error"))

	result, err := service.FindAll(ctx, filter)

	assert.Error(t, err)
	assert.Empty(t, result.Items)
	assert.Zero(t, result.Paging.TotalItems)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := users.NewUserService(mockRepo)
	ctx := newTestContext()

	objectId, _ := bson.ObjectIDFromHex("67f759abe02e4b2f3c2f69b6")
	expectedUsers := users.UserModel{
		ID:       objectId,
		Name:     "John Doe",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mockRepo.On("FindById", ctx, objectId.Hex()).
		Return(expectedUsers, nil)

	result, err := service.FindById(ctx, objectId.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUsers.Name, "John Doe")

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := users.NewUserService(mockRepo)
	ctx := newTestContext()

	objectId, _ := bson.ObjectIDFromHex("67f759abe02e4b2f3c2f69b6")

	mockRepo.On("FindById", ctx, objectId.Hex()).
		Return(users.UserModel{}, errors.New("user not found"))

	result, err := service.FindById(ctx, objectId.Hex())

	assert.Error(t, err)
	assert.Equal(t, result, users.UserModel{})

	mockRepo.AssertExpectations(t)
}
