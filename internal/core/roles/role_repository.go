package roles

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
	c := ctx.Request().Context()

	_, err := r.collection.InsertOne(c, role)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.NewConflict("data already exists")
		}
		return utils.NewInternal("failed to create data")
	}

	return nil
}

func (r *RoleRepository) FindById(ctx echo.Context, id string) (RoleModel, error) {
	c := ctx.Request().Context()

	var role RoleModel
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return RoleModel{}, utils.NewBadRequest("invalid id format")
	}

	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(c, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return RoleModel{}, utils.NewBadRequest("data not found")
		}

		return RoleModel{}, utils.NewInternal("failed to find data")
	}

	return role, nil
}

func (r *RoleRepository) FindAll(ctx echo.Context, filter *shared.PaginationFilter) ([]RoleModel, int, error) {
	c := ctx.Request().Context()

	var roles []RoleModel
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))
	// opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(c, bson.M{}, opts)
	if err != nil {
		return nil, 0, utils.NewInternal("failed to query data")
	}
	defer cursor.Close(c)

	if err := cursor.All(c, &roles); err != nil {
		return nil, 0, utils.NewInternal("failed to decode data")
	}

	totalItems, err = r.collection.CountDocuments(c, bson.M{})
	if err != nil {
		return nil, 0, utils.NewInternal("failed to count documents")
	}

	return roles, int(totalItems), nil
}

func (r *RoleRepository) Update(ctx echo.Context, id string, role *entities.Role) error {
	c := ctx.Request().Context()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}

	filter := bson.M{"_id": objectId}
	result := r.collection.FindOneAndUpdate(c, filter, bson.M{
		"$set": bson.M{
			"name":        role.Name,
			"permissions": role.Permissions,
			"updated_at":  time.Now(),
		}})

	if result.Err() != nil {
		return utils.NewInternal("failed to update data")
	}

	return nil
}

func (r *RoleRepository) Delete(ctx echo.Context, id string) error {
	c := ctx.Request().Context()

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}

	filter := bson.M{"_id": objectId}

	result := r.collection.FindOneAndDelete(c, filter)
	if result.Err() != nil {
		return utils.NewInternal("failed to delete data")
	}

	return nil
}

func (r *RoleRepository) AssignUser(ctx echo.Context, userId string, roleId string) error {
	c := ctx.Request().Context()

	userCollection := r.app.DB.Collection("users")
	objectUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}

	objectRoleID, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}

	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$addToSet": bson.M{"roles": objectRoleID},
	}

	_, err = userCollection.UpdateOne(c, filter, update)
	if err != nil {
		return utils.NewInternal("failed to assign role to user")
	}

	return nil
}

func (r *RoleRepository) UnassignUser(ctx echo.Context, userId string, roleId string) error {
	c := ctx.Request().Context()

	userCollection := r.app.DB.Collection("users")
	objectUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}

	objectRoleID, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return utils.NewBadRequest("invalid id format")
	}
	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$pull": bson.M{"roles": objectRoleID},
	}
	_, err = userCollection.UpdateOne(c, filter, update)
	if err != nil {
		return utils.NewInternal("failed to unassign role to user")
	}
	return nil
}
