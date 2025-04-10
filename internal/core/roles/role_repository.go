package roles

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

type RoleRepository struct {
	app        *app.Apps
	collection *mongo.Collection
}

func NewRoleRepository(app *app.Apps) *RoleRepository {
	return &RoleRepository{
		app:        app,
		collection: app.DB.Collection("roles"),
	}
}

func (r *RoleRepository) Create(ctx echo.Context, role *entities.Role) error {
	_, err := r.collection.InsertOne(ctx.Request().Context(), role)
	return err
}

func (r *RoleRepository) FindById(ctx echo.Context, id string) (RoleModel, error) {
	var role RoleModel
	objectID, _ := bson.ObjectIDFromHex(id)

	filter := bson.M{"_id": objectID}
	err := r.collection.FindOne(ctx.Request().Context(), filter).Decode(&role)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return RoleModel{}, nil
		}
		return RoleModel{}, err
	}
	return role, nil
}

func (r *RoleRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error) {
	var roles []RoleModel
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))
	// opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx.Request().Context(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx.Request().Context())

	if err := cursor.All(ctx.Request().Context(), &roles); err != nil {
		return nil, 0, err
	}

	totalItems, err = r.collection.CountDocuments(ctx.Request().Context(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return roles, int(totalItems), nil
}

func (r *RoleRepository) Update(ctx echo.Context, id string, role *entities.Role) error {
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	result := r.collection.FindOneAndUpdate(ctx.Request().Context(), filter, bson.M{
		"$set": bson.M{
			"name":        role.Name,
			"permissions": role.Permissions,
			"updated_at":  time.Now(),
		}})

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return result.Err()
	}
	return nil
}

func (r *RoleRepository) Delete(ctx echo.Context, id string) error {
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	result := r.collection.FindOneAndDelete(ctx.Request().Context(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return result.Err()
	}
	return nil
}

func (r *RoleRepository) AssignUser(ctx echo.Context, userId string, roleId string) error {
	userCollection := r.app.DB.Collection("users")
	objectUserID, _ := bson.ObjectIDFromHex(userId)
	objectRoleID, _ := bson.ObjectIDFromHex(roleId)

	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$addToSet": bson.M{"roles": objectRoleID},
	}
	_, err := userCollection.UpdateOne(ctx.Request().Context(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

func (r *RoleRepository) UnassignUser(ctx echo.Context, userId string, roleId string) error {
	userCollection := r.app.DB.Collection("users")
	objectUserID, _ := bson.ObjectIDFromHex(userId)
	objectRoleID, _ := bson.ObjectIDFromHex(roleId)

	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$pull": bson.M{"roles": objectRoleID},
	}
	_, err := userCollection.UpdateOne(ctx.Request().Context(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}
