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
	"net/url"
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
	// GitCode webhook 暂时不包含文件变更信息，返回空列表
	// TODO: 根据实际 GitCode webhook 格式更新
	return []string{}
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

// pipelineFromPullRequest extracts the Pipeline data from a GitCode pull_request hook.
func pipelineFromPullRequest(hook *pullRequestHook) *model.Pipeline {
	avatar := expandAvatar(
		hook.Repo.HTTPURLToRepo,
		fixMalformedAvatar(hook.Sender.AvatarURL),
	)

	event := model.EventPull
	if hook.Action == actionClose {
		event = model.EventPullClosed
	}

	pipeline := &model.Pipeline{
		Event:    event,
		Commit:   hook.PullRequest.Head.SHA,
		ForgeURL: fmt.Sprintf("%s/%s/pulls/%d", defaultURL, hook.Repo.FullName, hook.Number),
		Ref:      fmt.Sprintf("refs/pull/%d/head", hook.Number),
		Branch:   hook.PullRequest.Base.Ref,
		Message:  hook.PullRequest.Title,
		Author:   hook.Sender.Login,
		Avatar:   avatar,
		Sender:   hook.Sender.Login,
		Email:    hook.Sender.Email,
		Title:    hook.PullRequest.Title,
		Refspec: fmt.Sprintf("%s:%s",
			hook.PullRequest.Head.Ref,
			hook.PullRequest.Base.Ref,
		),
		PullRequestLabels: []string{}, // GitCode 暂时不支持标签
		FromFork:          false,      // GitCode 暂时不支持 fork 检测
	}

	return pipeline
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

// matchingHooks return matching hooks.
func matchingHooks(hooks []*Hook, rawURL string) *Hook {
	link, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	for _, hook := range hooks {
		hookURL, err := url.Parse(hook.URL)
		if err == nil && hookURL.Host == link.Host {
			return hook
		}
	}
	return nil
}
