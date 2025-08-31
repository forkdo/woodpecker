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
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

// TestGitCodeIntegration 测试 GitCode 集成功能
func TestGitCodeIntegration(t *testing.T) {
	// 跳过集成测试，除非设置了环境变量
	if os.Getenv("GITCODE_INTEGRATION_TEST") == "" {
		t.Skip("Skipping GitCode integration test. Set GITCODE_INTEGRATION_TEST=1 to run.")
	}

	token := os.Getenv("GITCODE_TOKEN")
	if token == "" {
		t.Skip("Skipping GitCode integration test. Set GITCODE_TOKEN to run.")
	}

	// 创建 GitCode forge 实例
	forge, err := New(Opts{
		URL:               "https://gitcode.com",
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		SkipVerify:        false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, forge)

	ctx := context.Background()

	// 测试认证
	t.Run("Auth", func(t *testing.T) {
		username, err := forge.Auth(ctx, token, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, username)
		t.Logf("Authenticated user: %s", username)
	})

	// 测试获取用户仓库
	t.Run("Repos", func(t *testing.T) {
		user := &model.User{
			AccessToken: token,
		}

		repos, err := forge.Repos(ctx, user)
		assert.NoError(t, err)
		assert.NotNil(t, repos)
		t.Logf("Found %d repositories", len(repos))

		if len(repos) > 0 {
			repo := repos[0]
			t.Logf("First repo: %s/%s", repo.Owner, repo.Name)
			assert.NotEmpty(t, repo.Owner)
			assert.NotEmpty(t, repo.Name)
			assert.NotEmpty(t, repo.FullName)
		}
	})

	// 测试获取分支
	t.Run("Branches", func(t *testing.T) {
		user := &model.User{
			AccessToken: token,
		}

		repos, err := forge.Repos(ctx, user)
		assert.NoError(t, err)

		if len(repos) > 0 {
			repo := repos[0]
			branches, err := forge.Branches(ctx, user, repo, &model.ListOptions{Page: 1, PerPage: 10})
			assert.NoError(t, err)
			assert.NotNil(t, branches)
			t.Logf("Found %d branches for repo %s", len(branches), repo.FullName)

			if len(branches) > 0 {
				t.Logf("First branch: %s", branches[0])
			}
		}
	})

	// 测试获取文件内容
	t.Run("File", func(t *testing.T) {
		user := &model.User{
			AccessToken: token,
		}

		repos, err := forge.Repos(ctx, user)
		assert.NoError(t, err)

		if len(repos) > 0 {
			repo := repos[0]
			pipeline := &model.Pipeline{
				Commit: repo.Branch, // 使用默认分支
			}

			// 尝试获取 README 文件
			content, err := forge.File(ctx, user, repo, pipeline, "README.md")
			if err == nil {
				assert.NotNil(t, content)
				t.Logf("README.md content length: %d bytes", len(content))
			} else {
				t.Logf("README.md not found: %v", err)
			}
		}
	})
}

// TestGitCodeForgeInterface 测试 GitCode 是否正确实现了 Forge 接口
func TestGitCodeForgeInterface(t *testing.T) {
	forge, err := New(Opts{
		URL:               "https://gitcode.com",
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		SkipVerify:        false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, forge)

	// 测试基本方法
	assert.Equal(t, "gitcode", forge.Name())
	assert.Equal(t, "https://gitcode.com", forge.URL())

	// 测试 Netrc 生成
	user := &model.User{
		Login:       "testuser",
		AccessToken: "test-token",
	}
	repo := &model.Repo{
		Clone: "https://gitcode.com/user/repo.git",
	}

	netrc, err := forge.Netrc(user, repo)
	assert.NoError(t, err)
	assert.NotNil(t, netrc)
	assert.Equal(t, "testuser", netrc.Login)
	assert.Equal(t, "test-token", netrc.Password)
	assert.Equal(t, "gitcode.com", netrc.Machine)
	assert.Equal(t, model.ForgeTypeGitCode, netrc.Type)
}
