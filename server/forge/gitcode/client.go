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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GitCodeClient GitCode API v5 客户端
type GitCodeClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewGitCodeClient 创建新的 GitCode 客户端
func NewGitCodeClient(token string, skipVerify bool) *GitCodeClient {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	return &GitCodeClient{
		baseURL:    defaultURL,
		token:      token,
		httpClient: httpClient,
	}
}

// makeRequest 发送 HTTP 请求
func (c *GitCodeClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// GitCode API 基础 URL，使用独立的 API 域名
	url := fmt.Sprintf("%s%s", defaultAPI, endpoint)

	// GitCode 使用 access_token 查询参数进行认证
	if c.token != "" {
		// 如果 URL 已经有查询参数，添加 &，否则添加 ?
		if strings.Contains(url, "?") {
			url += "&access_token=" + c.token
		} else {
			url += "?access_token=" + c.token
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置标准头部
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Woodpecker-CI")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// get 发送 GET 请求并解析 JSON 响应
func (c *GitCodeClient) get(ctx context.Context, endpoint string, result interface{}) error {
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("GitCode API Error %d for %s: %s", resp.StatusCode, endpoint, string(body))
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}

		log.Printf("GitCode API Response for %s: %s", endpoint, string(body))

		if err := json.Unmarshal(body, result); err != nil {
			log.Printf("GitCode API JSON decode error for %s: %v, body: %s", endpoint, err, string(body))
			return fmt.Errorf("decode JSON response: %w", err)
		}
	}

	return nil
}

// post 发送 POST 请求
func (c *GitCodeClient) post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	resp, err := c.makeRequest(ctx, "POST", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// delete 发送 DELETE 请求
func (c *GitCodeClient) delete(ctx context.Context, endpoint string) error {
	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GitCode API 数据结构

// User GitCode 用户信息
type User struct {
	ID        interface{} `json:"id"` // GitCode API 可能返回字符串或数字
	Login     string      `json:"login"`
	FullName  string      `json:"full_name"`
	Email     string      `json:"email"`
	AvatarURL string      `json:"avatar_url"`
}

// Repository GitCode 仓库信息 (基于实际 API 响应)
type Repository struct {
	// 基本信息
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	HumanName   string `json:"human_name"`
	URL         string `json:"url"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`

	// 命名空间
	Namespace struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		HTMLURL string `json:"html_url"`
	} `json:"namespace"`

	// URLs
	SSHURLToRepo  string `json:"ssh_url_to_repo"`
	HTTPURLToRepo string `json:"http_url_to_repo"`
	WebURL        string `json:"web_url"`
	ReadmeURL     string `json:"readme_url"`

	// 时间戳
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// 创建者
	Creator struct {
		ID       string `json:"id"`
		ArtsID   string `json:"arts_id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Photo    string `json:"photo"`
	} `json:"creator"`

	// 分支信息
	DefaultBranch string `json:"default_branch"`

	// Fork 信息
	Fork bool `json:"fork"` // 0 = original, 1 = fork

	// 所有者
	Owner struct {
		ID    string `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"owner"`

	// 分配者
	Assigner struct {
		ID    string `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"assigner"`

	// 可见性 (GitCode 使用 integer 类型)
	Private  bool `json:"private"`  // 0 = public, 1 = private
	Public   bool `json:"public"`   // 0 = private, 1 = public
	Internal bool `json:"internal"` // 0 = external, 1 = internal

	// 统计信息
	ForksCount      int `json:"forks_count"`
	StargazersCount int `json:"stargazers_count"`
	WatchersCount   int `json:"watchers_count"`
	OpenIssuesCount int `json:"open_issues_count"`
	AssigneesNumber int `json:"assignees_number"`

	// 企业信息
	Enterprise struct {
		ID      int    `json:"id"`
		Path    string `json:"path"`
		HTMLURL string `json:"html_url"`
		Type    string `json:"type"`
	} `json:"enterprise"`

	// 其他字段
	Members             []string      `json:"members"`
	ProjectLabels       []interface{} `json:"project_labels"`
	License             string        `json:"license"`
	IssueTemplateSource string        `json:"issue_template_source"`
}

// Branch GitCode 分支信息
type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

// PullRequest GitCode Pull Request 信息
type PullRequest struct {
	ID     int64  `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Head   struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
}

// Hook GitCode Webhook 信息
type Hook struct {
	ID     int64    `json:"id"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

// Tree GitCode 目录树结构 <yes>
type Tree struct {
	Tree []TreeEntry `json:"tree"`
}

// TreeEntry 目录树条目 <yes>
type TreeEntry struct {
	SHA  string `json:"sha"`
	Name string `json:"name"`
	Type string `json:"type"` // "blob", "tree"
	Path string `json:"path"`
	Mode string `json:"mode"`
	MD5  string `json:"md5"`
}

// CreateHookRequest 创建 Webhook 请求
type CreateHookRequest struct {
	URL         string   `json:"url"`
	ContentType string   `json:"content_type"`
	Events      []string `json:"events"`
	Active      bool     `json:"active"`
}

// GitCode API 方法

// GetUser 获取当前用户信息
func (c *GitCodeClient) GetUser(ctx context.Context) (*User, error) {
	var user User
	log.Printf("GitCode API: Getting user info with token: %s...", c.token[:10])
	err := c.get(ctx, "/user", &user)
	if err != nil {
		log.Printf("GitCode API Error getting user: %v", err)
		return nil, err
	}
	log.Printf("GitCode API Success: Got user %s", user.Login)
	return &user, err
}

// GetUserRepos 获取用户仓库列表
func (c *GitCodeClient) GetUserRepos(ctx context.Context, page, limit int) ([]*Repository, error) {
	endpoint := fmt.Sprintf("/user/repos?page=%d&per_page=%d&sort=updated&direction=desc", page, limit)
	log.Printf("Calling GitCode API endpoint: %s", endpoint)

	var repos []*Repository
	err := c.get(ctx, endpoint, &repos)
	if err != nil {
		log.Printf("GitCode API Error: %v", err)
		return nil, err
	}

	log.Printf("GitCode API Success: Got %d repos", len(repos))
	for i, repo := range repos {
		log.Printf("Repo %d: ID=%v, Name=%s, FullName=%s, Private=%v", i, repo.ID, repo.Name, repo.FullName, repo.Private)
	}
	return repos, nil
}

// GetRepo 获取仓库信息
func (c *GitCodeClient) GetRepo(ctx context.Context, owner, repo string) (*Repository, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s", owner, repo)
	var repository Repository
	err := c.get(ctx, endpoint, &repository)
	return &repository, err
}

// GetBranches 获取分支列表
func (c *GitCodeClient) GetBranches(ctx context.Context, owner, repo string) ([]*Branch, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/branches", owner, repo)
	var branches []*Branch
	err := c.get(ctx, endpoint, &branches)
	return branches, err
}

// GetBranch 获取分支信息
func (c *GitCodeClient) GetBranch(ctx context.Context, owner, repo, branch string) (*Branch, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/branches/%s", owner, repo, branch)
	var branchInfo Branch
	err := c.get(ctx, endpoint, &branchInfo)
	return &branchInfo, err
}

// GetPullRequests 获取 PR 列表
func (c *GitCodeClient) GetPullRequests(ctx context.Context, owner, repo string) ([]*PullRequest, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	var prs []*PullRequest
	err := c.get(ctx, endpoint, &prs)
	return prs, err
}

// GetFileContent 获取文件内容
func (c *GitCodeClient) GetFileContent(ctx context.Context, owner, repo, path, ref string) ([]byte, error) {
	// GitCode 使用不同的端点获取文件内容
	endpoint := fmt.Sprintf("/repos/%s/%s/raw/%s", owner, repo, path)
	if ref != "" {
		endpoint += "?ref=" + url.QueryEscape(ref)
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// GetTree 获取目录树结构 (基于 GitCode 文档)
func (c *GitCodeClient) GetTree(ctx context.Context, owner, repo, sha string, recursive bool) (*Tree, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/git/trees/%s", owner, repo, sha)
	if recursive {
		endpoint += "?recursive=1"
	}

	var tree Tree
	err := c.get(ctx, endpoint, &tree)
	return &tree, err
}

// CreateHook 创建 Webhook
func (c *GitCodeClient) CreateHook(ctx context.Context, owner, repo string, hook *CreateHookRequest) (*Hook, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/hooks", owner, repo)
	var result Hook
	err := c.post(ctx, endpoint, hook, &result)
	return &result, err
}

// GetHooks 获取 Webhook 列表
func (c *GitCodeClient) GetHooks(ctx context.Context, owner, repo string) ([]*Hook, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/hooks", owner, repo)
	var hooks []*Hook
	err := c.get(ctx, endpoint, &hooks)
	return hooks, err
}

// DeleteHook 删除 Webhook
func (c *GitCodeClient) DeleteHook(ctx context.Context, owner, repo string, hookID int64) error {
	endpoint := fmt.Sprintf("/repos/%s/%s/hooks/%d", owner, repo, hookID)
	return c.delete(ctx, endpoint)
}
