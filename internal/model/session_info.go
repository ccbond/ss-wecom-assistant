package model

import "gorm.io/gorm"

// SessionInfo is the model for session info.
type SessionInfo struct {
	gorm.Model
	UserID   string `gorm:"column:user_id;index"`
	NickName string `gorm:"column:nick_name"`
	Question string `gorm:"column:question"`
	Answer   string `gorm:"column:answer"`
}
