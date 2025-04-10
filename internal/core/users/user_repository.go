package users

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
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
	_, err := u.collection.InsertOne(ctx.Request().Context(), user)
	return err
}

func (u *UserRepository) FindByEmail(ctx echo.Context, email string) (UserModel, error) {
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

	cursor, err := u.collection.Aggregate(ctx.Request().Context(), pipeline)
	if err != nil {
		return UserModel{}, err
	}
	defer cursor.Close(ctx.Request().Context())

	var users []UserModel
	if err := cursor.All(ctx.Request().Context(), &users); err != nil {
		return UserModel{}, err
	}

	if len(users) == 0 {
		return UserModel{}, nil
	}

	return users[0], nil
}

func (u *UserRepository) FindById(ctx echo.Context, id string) (UserModel, error) {
	objectID, _ := bson.ObjectIDFromHex(id)
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

	cursor, err := u.collection.Aggregate(ctx.Request().Context(), pipeline)
	if err != nil {
		return UserModel{}, err
	}
	defer cursor.Close(ctx.Request().Context())

	var users []UserModel
	if err := cursor.All(ctx.Request().Context(), &users); err != nil {
		return UserModel{}, err
	}
	if len(users) == 0 {
		return UserModel{}, nil
	}
	return users[0], nil
}

func (u *UserRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]UserModelResponse, int, error) {
	var roles []UserModelResponse
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))
	// opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := u.collection.Find(ctx.Request().Context(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx.Request().Context())

	if err := cursor.All(ctx.Request().Context(), &roles); err != nil {
		return nil, 0, err
	}

	totalItems, err = u.collection.CountDocuments(ctx.Request().Context(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return roles, int(totalItems), nil

	// var users []entities.User
	// var totalItems int64

	// query := u.app.DB.WithContext(ctx.Request().Context())

	// // Hitung total data sebelum pagination
	// if err := query.Model(&entities.User{}).Count(&totalItems).Error; err != nil {
	// 	return nil, 0, err
	// }

	// // Query data dengan pagination
	// result := query.Scopes(utils.Paginate(filter)).
	// 	Select([]string{"id", "name", "email", "created_at"}).
	// 	Find(&users)
	// if result.Error != nil {
	// 	return nil, 0, result.Error
	// }

	// // Konversi ke response model
	// var userModels []UserModelResponse
	// for _, user := range users {
	// 	userModels = append(userModels, UserModelResponse{
	// 		ID:        (user.ID).String(),
	// 		Email:     user.Email,
	// 		Name:      user.Name,
	// 		CreatedAt: user.CreatedAt,
	// 	})
	// }

	// return userModels, int(totalItems), nil
}

func (u *UserRepository) Update(ctx echo.Context, id string, user *entities.User) error {
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	result := u.collection.FindOneAndUpdate(ctx.Request().Context(), filter, bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": time.Now(),
		}})

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return result.Err()
	}
	return nil
}

func (u *UserRepository) Delete(ctx echo.Context, id string) error {
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	result := u.collection.FindOneAndDelete(ctx.Request().Context(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return result.Err()
	}
	return nil
}
