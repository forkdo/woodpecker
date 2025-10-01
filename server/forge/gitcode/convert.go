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
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

// toRepo 将 GitCode Repository 转换为 Woodpecker Repo
func toRepo(from *Repository) *model.Repo {
	// GitCode 直接提供 full_name 字段
	fullName := from.FullName

	// 解析 owner 和 name
	var owner, name string
	if fullName != "" {
		parts := strings.Split(fullName, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			name = parts[1]
		}
	}

	// 备用方案：使用 namespace.path 或 name 字段
	if owner == "" && from.Namespace.Path != "" {
		owner = from.Namespace.Path
	}
	if name == "" && from.Name != "" {
		name = from.Name
	}

	// GitCode 使用 web_url 作为仓库页面 URL
	forgeURL := from.WebURL
	if forgeURL == "" {
		forgeURL = fmt.Sprintf("%s/%s", defaultURL, fullName)
	}

	// Clone URL 处理 - 优先使用 http_url_to_repo
	cloneURL := from.HTTPURLToRepo
	if cloneURL == "" {
		cloneURL = fmt.Sprintf("%s/%s.git", defaultURL, fullName)
	}

	// SSH URL
	cloneSSH := from.SSHURLToRepo
	if cloneSSH == "" {
		cloneSSH = fmt.Sprintf("git@%s:%s.git", strings.TrimPrefix(defaultURL, "https://"), fullName)
	}

	// 处理私有状态 - GitCode 使用 bool 类型
	isPrivate := from.Private

	// 处理权限 - 使用 GitCode API 返回的权限信息
	canPull := from.Permission.Pull
	canPush := from.Permission.Push
	canAdmin := from.Permission.Admin

	// 获取头像 - 使用创建者的头像
	avatar := ""
	if from.Creator.Photo != "" {
		avatar = expandAvatar(defaultURL, from.Creator.Photo)
	}

	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(fmt.Sprintf("%d", from.ID)),
		Owner:         owner,
		Name:          name,
		FullName:      fullName,
		Avatar:        avatar,
		ForgeURL:      forgeURL,
		Clone:         cloneURL,
		CloneSSH:      cloneSSH,
		Branch:        from.DefaultBranch,

		IsSCMPrivate: isPrivate,
		Perm: &model.Perm{
			Pull:  canPull,
			Push:  canPush,
			Admin: canAdmin,
		},
	}
}

// toTeam 将 GitCode Organization 转换为 Woodpecker Team
func toTeam(from *User, baseURL string) *model.Team {
	avatar := expandAvatar(baseURL, from.AvatarURL)
	return &model.Team{
		Login:  from.Login,
		Avatar: avatar,
	}
}

// expandAvatar 扩展头像 URL
func expandAvatar(baseURL, avatarURL string) string {
	if avatarURL == "" {
		return ""
	}
	if strings.HasPrefix(avatarURL, "http") {
		return avatarURL
	}
	if strings.HasPrefix(avatarURL, "//") {
		return "https:" + avatarURL
	}
	if strings.HasPrefix(avatarURL, "/") {
		u, _ := url.Parse(baseURL)
		u.Path = avatarURL
		return u.String()
	}
	return path.Join(baseURL, avatarURL)
}

// convertStatus 将 Woodpecker 状态转换为 GitCode 状态
func convertStatus(status model.StatusValue) string {
	switch status {
	case model.StatusPending, model.StatusBlocked:
		return "pending"
	case model.StatusRunning:
		return "running"
	case model.StatusSuccess:
		return "success"
	case model.StatusFailure:
		return "failure"
	case model.StatusKilled:
		return "cancelled"
	case model.StatusDeclined:
		return "cancelled"
	case model.StatusError:
		return "error"
	default:
		return "failure"
	}
}

// convertPullRequest 将 GitCode PullRequest 转换为 Woodpecker PullRequest
func convertPullRequest(from *PullRequest) *model.PullRequest {
	return &model.PullRequest{
		Index: model.ForgeRemoteID(strconv.Itoa(from.Number)),
		Title: from.Title,
	}
}

// convertBranch 将 GitCode Branch 转换为分支名
func convertBranch(from *Branch) string {
	return from.Name
}

// convertCommit 将 GitCode Commit 转换为 Woodpecker Commit
func convertCommit(sha, url string) *model.Commit {
	return &model.Commit{
		SHA:      sha,
		ForgeURL: url,
	}
}
