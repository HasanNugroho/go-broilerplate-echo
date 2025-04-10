package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email     string          `bson:"email" json:"email"`
	Name      string          `bson:"name" json:"name"`
	Password  string          `bson:"password" json:"password"`
	Roles     []bson.ObjectID `bson:"roles" json:"roles"`
	CreatedAt time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time       `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
