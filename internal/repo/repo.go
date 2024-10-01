// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package repo

import (
	"ss-wecom-assistant/internal/model"
	"time"
)

// SessionInfo is the interface for session info.
type SessionInfo interface {
	Create(*model.SessionInfo) error
	Get(uint64) (*model.SessionInfo, error)
	Update(*model.SessionInfo) error
	Delete(string) error
	List([]uint64) ([]*model.SessionInfo, error)
	ListRecentlyByUserID(string, int) ([]*model.SessionInfo, error)
	GetLatestByUserID(string) (*model.SessionInfo, error)
	ListByID(uint64, int) ([]*model.SessionInfo, error)
	ListByTimeInterval(time.Time, time.Time) (int, []*model.SessionInfo, error)
}

type User interface{}
