package entity

import (
	"time"
)

type User struct {
	ID               uint64     `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned"`
	Name             string     `gorm:"column:name;type:varchar(255);not null"`
	Email            string     `gorm:"column:email;type:varchar(255);not null"`
	Password         string     `gorm:"column:password;type:varchar(255);not null"`
	VerificationCode *string    `gorm:"column:verification_code;type:varchar(255);null"`
	VerifiedAt       *time.Time `gorm:"column:verified_at;type:datetime;null"`
	CreatedAt        time.Time  `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp on update current_timestamp"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;type:datetime;index"`
}

func (User) TableName() string { return "users" }
