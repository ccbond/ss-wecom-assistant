// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package repo

import (
	"errors"
	"ss-wecom-assistant/internal/model"
	"time"

	"gorm.io/gorm"
)

type sessionInfoDao struct {
	*gorm.DB
}

// NewSessionInfo is creating session info.
func NewSessionInfo(db *gorm.DB) SessionInfo {
	return &sessionInfoDao{db}
}

// Create is creating session info.
func (dao *sessionInfoDao) Create(session *model.SessionInfo) error {
	if err := dao.DB.Create(session).Error; err != nil {
		return err
	}
	return nil
}

// Get is getting session info.
func (dao *sessionInfoDao) Get(ID uint64) (*model.SessionInfo, error) {
	var session model.SessionInfo
	if err := dao.DB.Where("id = ?", ID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("session info not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetBySessionID is getting session info.
func (dao *sessionInfoDao) Update(session *model.SessionInfo) error {
	if err := dao.DB.Save(session).Error; err != nil {
		return err
	}
	return nil
}

// Delete is deleting session info.
func (dao *sessionInfoDao) Delete(sessionID string) error {
	if err := dao.DB.Where("id = ?", sessionID).Delete(model.SessionInfo{}).Error; err != nil {
		return err
	}
	return nil
}

// List is listing session info.
func (dao *sessionInfoDao) List(IDs []uint64) ([]*model.SessionInfo, error) {
	var sessions []*model.SessionInfo
	if err := dao.DB.Where("id IN (?)", IDs).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// ListRecentlyByUserID is listing session info.
func (dao *sessionInfoDao) ListRecentlyByUserID(userID string, limit int) ([]*model.SessionInfo, error) {
	var sessions []*model.SessionInfo
	if err := dao.DB.Where("user_id = ?", userID).Order("id DESC").Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetLatestByUserID is getting the latest session info.
func (dao *sessionInfoDao) GetLatestByUserID(userID string) (*model.SessionInfo, error) {
	var session model.SessionInfo
	if err := dao.DB.Where("user_id = ?", userID).Order("created_at desc").First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle no record found
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (dao *sessionInfoDao) ListByID(startID uint64, limit int) ([]*model.SessionInfo, error) {
	var sessions []*model.SessionInfo
	if err := dao.DB.Where("id >=?", startID).Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (dao *sessionInfoDao) ListByTimeInterval(startTime, endTime time.Time) (int, []*model.SessionInfo, error) {
	var sessions []*model.SessionInfo
	if err := dao.DB.Where("created_at >= ? AND created_at <= ?", startTime, endTime).Find(&sessions).Error; err != nil {
		return 0, nil, err
	}
	return len(sessions), sessions, nil
}
