// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID         string `gorm:"column:user_id"`
	UnionID        string `gorm:"column:union_id"`
	OpenID         string `gorm:"column:open_id;unique"`
	Email          string `gorm:"column:email;default:null"`
	Phone          string `gorm:"column:phone;default:null"`
	SubscriptionID uint   `gorm:"column:subscription_id;default:null"`
}
