package entity

import (
	"time"

	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
)

type User struct {
	ID        utils.ULID `gorm:"column:id;primaryKey;index"`
	Email     string     `gorm:"column:email"`
	Name      string     `gorm:"column:name"`
	Password  string     `gorm:"column:password"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at:autoUpdateTime"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
