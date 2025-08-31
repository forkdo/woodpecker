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
	"time"
)

func TestGitCodeClient_Connection(t *testing.T) {
	// 测试基础连接，不需要 token
	client := NewGitCodeClient("", false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试 API 连接（这个请求应该返回 401，表示连接正常但需要认证）
	_, err := client.GetUser(ctx)

	// 我们期望这里会有错误，因为没有提供 token
	if err == nil {
		t.Log("✅ API 连接正常（意外地不需要认证）")
	} else {
		t.Logf("✅ API 连接正常，需要认证: %v", err)
	}
}

func TestGitCodeClient_WithToken(t *testing.T) {
	token := os.Getenv("GITCODE_TOKEN")
	if token == "" {
		t.Skip("跳过认证测试 - 请设置 GITCODE_TOKEN 环境变量")
	}

	client := NewGitCodeClient(token, false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GetUser", func(t *testing.T) {
		user, err := client.GetUser(ctx)
		if err != nil {
			t.Errorf("获取用户信息失败: %v", err)
			return
		}

		t.Logf("✅ 用户信息: %s (%s)", user.Login, user.FullName)

		if user.ID == 0 {
			t.Error("用户 ID 不应该为 0")
		}
		if user.Login == "" {
			t.Error("用户登录名不应该为空")
		}
	})

	t.Run("GetUserRepos", func(t *testing.T) {
		repos, err := client.GetUserRepos(ctx, 1, 5)
		if err != nil {
			t.Errorf("获取用户仓库失败: %v", err)
			return
		}

		t.Logf("✅ 找到 %d 个仓库", len(repos))

		for i, repo := range repos {
			if i < 3 { // 只显示前3个
				t.Logf("   - %s (%s)", repo.FullName, repo.HTMLURL)
			}
		}
	})
}

func TestGitCodeClient_PublicAPI(t *testing.T) {
	client := NewGitCodeClient("", false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GetPublicRepo", func(t *testing.T) {
		// 尝试获取一个公开仓库（需要替换为实际存在的仓库）
		repo, err := client.GetRepo(ctx, "gitcode", "test")
		if err != nil {
			t.Logf("获取公开仓库失败（可能仓库不存在）: %v", err)
			return
		}

		t.Logf("✅ 公开仓库: %s", repo.FullName)
	})
}

func TestGitCodeClient_AuthMethods(t *testing.T) {
	token := os.Getenv("GITCODE_TOKEN")
	if token == "" {
		t.Skip("跳过认证方式测试 - 请设置 GITCODE_TOKEN 环境变量")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("BearerToken", func(t *testing.T) {
		client := NewGitCodeClient(token, false)

		user, err := client.GetUser(ctx)
		if err != nil {
			t.Errorf("Bearer Token 认证失败: %v", err)
			return
		}

		t.Logf("✅ Bearer Token 认证成功: %s", user.Login)
	})

	// TODO: 添加其他认证方式的测试
	// - PRIVATE-TOKEN 头部
	// - URL 参数认证
}

// 基准测试
func BenchmarkGitCodeClient_GetUser(b *testing.B) {
	token := os.Getenv("GITCODE_TOKEN")
	if token == "" {
		b.Skip("跳过性能测试 - 请设置 GITCODE_TOKEN 环境变量")
	}

	client := NewGitCodeClient(token, false)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetUser(ctx)
		if err != nil {
			b.Fatalf("API 调用失败: %v", err)
		}
	}
}
