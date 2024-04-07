// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package tiktoken

import "github.com/pkoukk/tiktoken-go"

// CountTokenNums count token numbers
// Default model is gpt-3.5-turbo
func CountTokenNums(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		return 0, err
	}
	token := tkm.Encode(text, nil, nil)

	num_token := len(token)

	return num_token, nil
}

func TruncateToTokenLimit(text string, limit int) (string, int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		return "", 0, err
	}
	token := tkm.Encode(text, nil, nil)
	resultLength := len(text)
	if len(token) > limit {
		token = token[:limit]
		resultLength = limit
	}

	return tkm.Decode(token), resultLength, nil
}
