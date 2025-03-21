package entity

import (
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey;index"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at:autoUpdateTime"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
