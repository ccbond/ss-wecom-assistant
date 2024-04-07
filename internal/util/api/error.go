// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package api

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int
	Meta       map[string]interface{}
	RawBody    []byte
}

func (err *APIError) Error() string {
	if err.Meta != nil {
		metaBytes, _ := json.Marshal(err.Meta)
		return fmt.Sprintf("APIError: code=%d meta=%s", err.StatusCode, string(metaBytes))
	}
	return fmt.Sprintf("APIError: code=%d body=%s", err.StatusCode, string(err.RawBody))
}
