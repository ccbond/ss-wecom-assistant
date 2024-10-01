// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	OpenID        string `gorm:"column:open_id;unique"`
	NickName      string `gorm:"column:nick_name"`
	TotalMessages int    `gorm:"column:total_messages"`
}
