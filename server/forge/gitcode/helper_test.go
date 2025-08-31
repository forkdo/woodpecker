// Copyright 2024 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixMalformedAvatar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "triple slash",
			input:    "https://gitcode.com///avatars/1",
			expected: "//avatars/1",
		},
		{
			name:     "double avatars slash",
			input:    "https://gitcode.com//avatars/1",
			expected: "https://gitcode.com/avatars/1",
		},
		{
			name:     "normal url",
			input:    "https://gitcode.com/avatars/1",
			expected: "https://gitcode.com/avatars/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixMalformedAvatar(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandAvatar(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		rawURL   string
		expected string
	}{
		{
			name:     "absolute url",
			repo:     "https://gitcode.com/user/repo",
			rawURL:   "https://gitcode.com/avatars/1",
			expected: "https://gitcode.com/avatars/1",
		},
		{
			name:     "relative url",
			repo:     "https://gitcode.com/user/repo",
			rawURL:   "/avatars/1",
			expected: "https://gitcode.com/avatars/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandAvatar(tt.repo, tt.rawURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestConvertLabels 测试标签转换（GitCode 暂时不支持标签）
func TestConvertLabels(t *testing.T) {
	// GitCode 暂时不支持标签功能，跳过此测试
	t.Skip("GitCode labels not implemented yet")
}
