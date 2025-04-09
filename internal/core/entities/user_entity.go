package entities

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"gorm.io/gorm"
)

type User struct {
	ID        utils.ULID `gorm:"column:id;primaryKey;type:varchar(26)"`
	Email     string     `gorm:"column:email"`
	Name      string     `gorm:"column:name"`
	Password  string     `gorm:"column:password"`
	Roles     []Role     `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.NewULID()
	return nil
}
