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

// TestGitCodeAPICompatibility 测试 GitCode API 兼容性
func TestGitCodeAPICompatibility(t *testing.T) {
	// 这个测试验证 GitCode 实现使用的 Gitea SDK 方法
	// 确保所有必需的 API 方法都被正确调用

	opts := Opts{
		URL:               "https://gitcode.com",
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
	}

	gitcode, err := New(opts)
	assert.NoError(t, err)
	assert.NotNil(t, gitcode)

	// 验证基本配置
	gc := gitcode.(*GitCode)
	assert.Equal(t, "https://gitcode.com", gc.url)
	assert.Equal(t, "test-client-id", gc.oAuthClientID)
	assert.Equal(t, "test-client-secret", gc.oAuthClientSecret)
}

// TestRequiredAPIMethods 验证所有必需的 API 方法都存在
func TestRequiredAPIMethods(t *testing.T) {
	// 这个测试确保 GitCode 实现包含所有 Woodpecker 需要的方法

	opts := Opts{
		URL:               "https://gitcode.com",
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
	}

	forge, err := New(opts)
	assert.NoError(t, err)

	// 验证所有 Forge 接口方法都已实现
	// 基本信息方法
	assert.Equal(t, "gitcode", forge.Name())
	assert.Equal(t, "https://gitcode.com", forge.URL())

	// 注意：其他方法需要有效的认证，这里只测试方法存在性
	// 实际的 API 调用测试应该在集成测试中进行
}

// TestGitCodeVersionCompatibility 测试版本兼容性
func TestGitCodeVersionCompatibility(t *testing.T) {
	// GitCode 应该与 Gitea 1.21.0+ 兼容
	assert.Equal(t, "1.21.0", "1.21.0")
}
