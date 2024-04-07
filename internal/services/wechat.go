// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package services

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

type wechatService struct {
	weCom *work.Work
}

func NewWeComService(weCom *work.Work) WechatService {
	return &wechatService{
		weCom: weCom,
	}
}

func (w *wechatService) GetAccessToken() (string, error) {
	return "", nil
}
