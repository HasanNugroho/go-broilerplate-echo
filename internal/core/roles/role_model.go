package roles

import "go.mongodb.org/mongo-driver/v2/bson"

type RoleModel struct {
	ID          bson.ObjectID `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Permissions []string      `bson:"permissions" json:"permission"`
}
type RoleUpdateModel struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permission"`
}

type AssignRoleModel struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}
