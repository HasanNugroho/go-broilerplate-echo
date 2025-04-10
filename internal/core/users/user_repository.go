package users

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func convertBsonAToStringSlice(bsonArray bson.A) []string {
	var result []string
	for _, item := range bsonArray {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

type UserRepository struct {
	app        *app.Apps
	collection *mongo.Collection
}

func NewUserRepository(app *app.Apps) *UserRepository {
	return &UserRepository{
		app:        app,
		collection: app.DB.Collection("users"),
	}
}

func (u *UserRepository) Create(ctx echo.Context, user *entities.User) error {
	c := ctx.Request().Context()

	_, err := u.collection.InsertOne(c, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.NewConflict("email already exists")
		}
		return utils.NewInternal("failed to create user")
	}
	return nil
}

func (u *UserRepository) FindByEmail(ctx echo.Context, email string) (UserModel, error) {
	c := ctx.Request().Context()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "email", Value: email}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "roles"},
				{Key: "localField", Value: "roles"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "roles_data"},
			},
		}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := u.collection.Aggregate(c, pipeline)
	if err != nil {
		return UserModel{}, utils.NewInternal("failed to query data")
	}
	defer cursor.Close(c)

	var users []UserModel
	if err := cursor.All(c, &users); err != nil {
		return UserModel{}, utils.NewInternal("failed to decode data")
	}

	if len(users) == 0 {
		return UserModel{}, utils.NewNotFound("data not found")
	}

	return users[0], nil
}

func (u *UserRepository) FindById(ctx echo.Context, id string) (UserModel, error) {
	c := ctx.Request().Context()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return UserModel{}, utils.NewBadRequest("invalid user id")
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "roles"},
				{Key: "localField", Value: "roles"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "roles_data"},
			},
		}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := u.collection.Aggregate(c, pipeline)
	if err != nil {
		return UserModel{}, utils.NewInternal("failed to query data")
	}
	defer cursor.Close(c)

	var users []UserModel
	if err := cursor.All(c, &users); err != nil {
		return UserModel{}, utils.NewInternal("failed to decode data")
	}

	if len(users) == 0 {
		return UserModel{}, utils.NewNotFound("data not found")
	}

	return users[0], nil
}

func (u *UserRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error) {
	c := ctx.Request().Context()
	var roles []UserModelResponse
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))
	// opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := u.collection.Find(c, bson.M{}, opts)
	if err != nil {
		return []UserModelResponse{}, 0, utils.NewInternal("failed to query data")
	}
	defer cursor.Close(c)

	if err := cursor.All(c, &roles); err != nil {
		return []UserModelResponse{}, 0, utils.NewInternal("failed decode data")
	}

	totalItems, err = u.collection.CountDocuments(c, bson.M{})
	if err != nil {
		return []UserModelResponse{}, 0, utils.NewInternal("failed count user")
	}

	return roles, int(totalItems), nil
}

func (u *UserRepository) Update(ctx echo.Context, id string, user *entities.User) error {
	c := ctx.Request().Context()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.NewBadRequest("invalid user id")
	}

	filter := bson.M{"_id": objectId}
	err = u.collection.FindOneAndUpdate(c, filter, bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": time.Now(),
		}}).Err()

	if err != nil {
		return utils.NewInternal("failed to update user")
	}

	return nil
}

func (u *UserRepository) Delete(ctx echo.Context, id string) error {
	c := ctx.Request().Context()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.NewBadRequest("invalid user id")
	}

	filter := bson.M{"_id": objectId}

	result := u.collection.FindOneAndDelete(c, filter)
	if result.Err() != nil {
		return utils.NewInternal("failed to update user")
	}

	return nil
}
