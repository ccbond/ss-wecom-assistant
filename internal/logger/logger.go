// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package logger

import "github.com/SyntSugar/ss-infra-go/log"

var _logger *log.Logger

// Init initialize logger
func Init(loglevel string) error {
	var err error
	_logger, err = log.NewLogger(loglevel, "json", log.DefaultSamplingConfig(), "")
	return err
}

// Get get logger
func Get() *log.Logger {
	if _logger != nil {
		return _logger
	}
	return log.GlobalLogger()
}
