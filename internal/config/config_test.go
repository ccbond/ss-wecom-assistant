// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitSensitiveWords(t *testing.T) {
	testConfigFilePath := "../../conf/config.toml"
	Env = "local"

	Init(testConfigFilePath)
	config := Get()

	t.Logf("%v", config)
	require.NotNil(t, config)
}
