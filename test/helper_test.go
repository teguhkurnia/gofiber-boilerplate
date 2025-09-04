package test

import (
	"gofiber-boilerplate/internal/entity"
)

func ClearAll() {
	ClearUsersTable()
}

func ClearUsersTable() {
	err := db.Where("1 = 1").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed to clear users table: %v", err)
	}
}
