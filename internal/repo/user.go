package repo

import (
	"gorm.io/gorm"
)

type userRepoDao struct {
	*gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &userRepoDao{db}
}
