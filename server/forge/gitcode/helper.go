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
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

// toRepo converts a GitCode repository to a Woodpecker repository.
// 这个函数在 convert.go 中已经定义，这里移除重复定义

// toPerm converts a GitCode permission to a Woodpecker permission.

// toTeam converts a GitCode team to a Woodpecker team.
// 这个函数在 convert.go 中已经定义，这里移除重复定义

// pipelineFromPush extracts the Pipeline data from a GitCode push hook.
func pipelineFromPush(hook *pushHook) *model.Pipeline {
	// 使用用户头像
	avatar := expandAvatar(defaultURL, fixMalformedAvatar(hook.UserAvatar))

	var message string
	var link string

	// 使用提交信息
	if len(hook.Commits) > 0 {
		message = hook.Commits[0].Message
		link = hook.Commits[0].URL
	} else {
		message = hook.Message
		link = hook.Project.WebURL
	}

	return &model.Pipeline{
		Event:        model.EventPush,
		Commit:       hook.After,
		Ref:          hook.Ref,
		ForgeURL:     link,
		Branch:       strings.TrimPrefix(hook.Ref, "refs/heads/"),
		Message:      message,
		Avatar:       avatar,
		Author:       hook.UserUsername,
		Email:        hook.UserEmail,
		Timestamp:    time.Now().UTC().Unix(),
		Sender:       hook.UserUsername,
		ChangedFiles: getChangedFilesFromPushHook(hook),
	}
}

func getChangedFilesFromPushHook(hook *pushHook) []string {
	var files []string
	for _, commit := range hook.Commits {
		files = append(files, commit.Added...)
		files = append(files, commit.Modified...)
		files = append(files, commit.Removed...)
	}

	// 去重
	seen := make(map[string]bool)
	var uniqueFiles []string
	for _, file := range files {
		if !seen[file] {
			seen[file] = true
			uniqueFiles = append(uniqueFiles, file)
		}
	}

	return uniqueFiles
}

// pipelineFromTag extracts the Pipeline data from a GitCode tag hook.
func pipelineFromTag(hook *pushHook) *model.Pipeline {
	avatar := expandAvatar(defaultURL, fixMalformedAvatar(hook.UserAvatar))
	ref := strings.TrimPrefix(hook.Ref, "refs/tags/")

	return &model.Pipeline{
		Event:     model.EventTag,
		Commit:    hook.After,
		Ref:       fmt.Sprintf("refs/tags/%s", ref),
		ForgeURL:  fmt.Sprintf("%s/src/tag/%s", hook.Project.WebURL, ref),
		Message:   fmt.Sprintf("created tag %s", ref),
		Avatar:    avatar,
		Author:    hook.UserUsername,
		Sender:    hook.UserUsername,
		Email:     hook.UserEmail,
		Timestamp: time.Now().UTC().Unix(),
	}
}

// pipelineFromPullRequestHook extracts the Pipeline data from a GitCode pull_request hook.
func pipelineFromPullRequestHook(hook *pullRequestHook) *model.Pipeline {
	avatar := expandAvatar(
		hook.Project.WebURL,
		fixMalformedAvatar(hook.User.AvatarURL),
	)

	event := model.EventPull
	if hook.MergeRequest.Action == "close" || hook.MergeRequest.State == "closed" {
		event = model.EventPullClosed
	}

	pipeline := &model.Pipeline{
		Event:    event,
		Commit:   hook.MergeRequest.LastCommit.ID,
		ForgeURL: hook.MergeRequest.URL,
		Ref:      fmt.Sprintf("refs/pull/%d/head", hook.MergeRequest.IID),
		Branch:   hook.MergeRequest.TargetBranch,
		Message:  hook.MergeRequest.Title,
		Author:   hook.User.Username,
		Avatar:   avatar,
		Sender:   hook.User.Username,
		Email:    hook.User.Email,
		Title:    hook.MergeRequest.Title,
		Refspec: fmt.Sprintf("%s:%s",
			hook.MergeRequest.SourceBranch,
			hook.MergeRequest.TargetBranch,
		),
		PullRequestLabels: []string{}, // GitCode 暂时不支持标签
		FromFork:          hook.MergeRequest.Source.ID != hook.MergeRequest.Target.ID,
	}

	return pipeline
}

// repoFromPullRequestHook extracts the Repository data from a GitCode pull_request hook.
func repoFromPullRequestHook(hook *pullRequestHook) *model.Repo {
	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(fmt.Sprintf("%d", hook.Project.ID)),
		Owner:         hook.Project.Namespace,
		Name:          hook.Project.Name,
		FullName:      hook.Project.PathWithNamespace,
		Avatar:        hook.Project.AvatarURL,
		ForgeURL:      hook.Project.WebURL,
		Clone:         hook.Project.GitHTTPURL,
		CloneSSH:      hook.Project.GitSSHURL,
		Branch:        hook.Project.DefaultBranch,
		IsSCMPrivate:  hook.Project.VisibilityLevel == 0,
		Perm: &model.Perm{
			Pull:  true,
			Push:  true,
			Admin: false,
		},
	}
}

func pipelineFromRelease(hook *releaseHook) *model.Pipeline {
	avatar := expandAvatar(
		hook.Repo.HTTPURLToRepo,
		fixMalformedAvatar(hook.Sender.AvatarURL),
	)

	return &model.Pipeline{
		Event:        model.EventRelease,
		Ref:          fmt.Sprintf("refs/tags/%s", hook.Release.TagName),
		ForgeURL:     fmt.Sprintf("%s/%s/releases/tag/%s", defaultURL, hook.Repo.FullName, hook.Release.TagName),
		Branch:       hook.Repo.DefaultBranch,
		Message:      fmt.Sprintf("created release %s", hook.Release.Name),
		Avatar:       avatar,
		Author:       hook.Sender.Login,
		Sender:       hook.Sender.Login,
		Email:        hook.Sender.Email,
		IsPrerelease: hook.Release.Prerelease,
	}
}

// parsePush parses a push hook from a read closer.
func parsePush(r io.Reader) (*pushHook, error) {
	push := new(pushHook)
	err := json.NewDecoder(r).Decode(push)
	return push, err
}

func parsePullRequest(r io.Reader) (*pullRequestHook, error) {
	pr := new(pullRequestHook)
	err := json.NewDecoder(r).Decode(pr)
	return pr, err
}

func parseRelease(r io.Reader) (*releaseHook, error) {
	pr := new(releaseHook)
	err := json.NewDecoder(r).Decode(pr)
	return pr, err
}

// fixMalformedAvatar fixes an avatar url if malformed (currently a known bug with gitcode).
func fixMalformedAvatar(url string) string {
	index := strings.Index(url, "///")
	if index != -1 {
		return url[index+1:]
	}
	index = strings.Index(url, "//avatars/")
	if index != -1 {
		return strings.ReplaceAll(url, "//avatars/", "/avatars/")
	}
	return url
}
