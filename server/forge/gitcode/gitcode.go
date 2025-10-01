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
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/rs/zerolog/log"
	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge/common"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	shared_utils "go.woodpecker-ci.org/woodpecker/v3/shared/utils"
)

const (
	defaultURL = "https://gitcode.com"            // Default GitCode URL
	defaultAPI = "https://api.gitcode.com/api/v5" // Default GitCode API URL

	// OAuth 端点
	authorizeTokenURL = "%s/oauth/authorize"
	accessTokenURL    = "%s/oauth/token"

	// API 配置
	defaultPageSize = 50
)

type Opts struct {
	OAuthClientID     string
	OAuthClientSecret string
}

type GitCode struct {
	oAuthClientID     string
	oAuthClientSecret string
	pageSize          int
}

func New(opts Opts) (forge.Forge, error) {
	return &GitCode{
		oAuthClientID:     opts.OAuthClientID,
		oAuthClientSecret: opts.OAuthClientSecret,
	}, nil
}

func (c *GitCode) Name() string {
	return "gitcode"
}

func (c *GitCode) URL() string {
	return defaultURL
}

func (c *GitCode) oauth2Config(ctx context.Context) (*oauth2.Config, context.Context) {
	return &oauth2.Config{
			ClientID:     c.oAuthClientID,
			ClientSecret: c.oAuthClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf(authorizeTokenURL, defaultURL),
				TokenURL: fmt.Sprintf(accessTokenURL, defaultURL),
			},
			RedirectURL: fmt.Sprintf("%s/authorize", server.Config.Server.OAuthHost),
		},

		context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			Proxy:           http.ProxyFromEnvironment,
		}})
}

func (c *GitCode) Login(ctx context.Context, req *forge_types.OAuthRequest) (*model.User, string, error) {
	config, oauth2Ctx := c.oauth2Config(ctx)
	redirectURL := config.AuthCodeURL(req.State)

	log.Debug().Msgf("GitCode OAuth config - AuthURL: %s, TokenURL: %s", config.Endpoint.AuthURL, config.Endpoint.TokenURL)

	// check the OAuth code
	if len(req.Code) == 0 {
		return nil, redirectURL, nil
	}

	log.Debug().Msgf("GitCode OAuth - Exchanging code: %s", req.Code)
	token, err := config.Exchange(oauth2Ctx, req.Code)
	if err != nil {
		log.Error().Err(err).Msgf("GitCode OAuth token exchange failed")
		return nil, redirectURL, err
	}

	client := NewGitCodeClient(token.AccessToken, false)
	account, err := client.GetUser(ctx)
	if err != nil {
		return nil, redirectURL, err
	}

	return &model.User{
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		Expiry:        token.Expiry.UTC().Unix(),
		Login:         account.Login,
		Email:         account.Email,
		ForgeRemoteID: model.ForgeRemoteID(fmt.Sprint(account.ID)),
		Avatar:        expandAvatar(defaultURL, account.AvatarURL),
	}, redirectURL, nil
}

func (c *GitCode) Auth(ctx context.Context, token, _ string) (string, error) {
	client := NewGitCodeClient(token, false)
	user, err := client.GetUser(ctx)
	if err != nil {
		return "", err
	}
	return user.Login, nil
}

func (c *GitCode) Refresh(ctx context.Context, user *model.User) (bool, error) {
	config, oauth2Ctx := c.oauth2Config(ctx)
	config.RedirectURL = ""

	source := config.TokenSource(oauth2Ctx, &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       time.Unix(user.Expiry, 0),
	})

	token, err := source.Token()
	if err != nil || len(token.AccessToken) == 0 {
		return false, err
	}

	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.Expiry = token.Expiry.UTC().Unix()
	return true, nil
}

func (c *GitCode) Teams(ctx context.Context, u *model.User) ([]*model.Team, error) {
	// GitCode 暂时不支持组织列表，返回空列表
	// TODO: 实现 GitCode 组织 API 支持
	return []*model.Team{}, nil
}

func (c *GitCode) TeamPerm(_ *model.User, _ string) (*model.Perm, error) {
	return nil, nil
}

func (c *GitCode) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	client := NewGitCodeClient(u.AccessToken, false)

	if remoteID.IsValid() {
		// GitCode 不支持直接通过 ID 获取仓库，需要从用户仓库列表中查找
		repos, err := shared_utils.Paginate(func(page int) ([]*Repository, error) {
			repos, err := client.GetUserRepos(ctx, page, defaultPageSize)
			return repos, err
		}, -1)
		if err != nil {
			return nil, err
		}

		// 查找匹配的仓库
		targetID := string(remoteID)
		for _, repo := range repos {
			if strconv.FormatInt(repo.ID, 10) == targetID {
				return toRepo(repo), nil
			}
		}
		return nil, fmt.Errorf("repository with ID %s not found", targetID)
	}

	// 通过 owner/name 获取仓库信息
	repo, err := client.GetRepo(ctx, owner, name)
	if err != nil {
		return nil, err
	}
	return toRepo(repo), nil
}

func (c *GitCode) Repos(ctx context.Context, u *model.User) ([]*model.Repo, error) {
	client := NewGitCodeClient(u.AccessToken, false)

	log.Debug().Msgf("GitCode: Getting repos for user %s", u.Login)

	repos, err := shared_utils.Paginate(func(page int) ([]*Repository, error) {
		repos, err := client.GetUserRepos(ctx, page, defaultPageSize)
		return repos, err
	}, -1)

	if err != nil {
		log.Error().Err(err).Msgf("GitCode: Failed to get repos for user %s", u.Login)
		return nil, err
	}

	log.Debug().Msgf("GitCode: Got %d repos from API for user %s", len(repos), u.Login)

	result := make([]*model.Repo, 0, len(repos))
	for i, repo := range repos {
		convertedRepo := toRepo(repo)
		log.Debug().Msgf("GitCode: Repo %d - Original: ID=%v, Name=%s, FullName=%s", i, repo.ID, repo.Name, repo.FullName)
		log.Debug().Msgf("GitCode: Repo %d - Converted: Owner=%s, Name=%s, FullName=%s, ForgeURL=%s", i, convertedRepo.Owner, convertedRepo.Name, convertedRepo.FullName, convertedRepo.ForgeURL)
		result = append(result, convertedRepo)
	}

	log.Debug().Msgf("GitCode: Returning %d converted repos for user %s", len(result), u.Login)
	return result, err
}

func (c *GitCode) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]byte, error) {
	client := NewGitCodeClient(u.AccessToken, false)

	// 确定要使用的 commit SHA 或分支名
	ref := b.Commit
	if ref == "" {
		// 如果没有指定 commit，使用默认分支
		ref = r.Branch
		if ref == "" {
			ref = "main" // 默认分支
		}
	}

	cfg, err := client.GetFileContent(ctx, r.Owner, r.Name, f, ref)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, errors.Join(err, &forge_types.ErrConfigNotFound{Configs: []string{f}})
		}
		return nil, err
	}
	return cfg, nil
}

func (c *GitCode) Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]*forge_types.FileMeta, error) {
	client := NewGitCodeClient(u.AccessToken, false)

	// 确定要使用的 commit SHA
	commitSHA := b.Commit
	if commitSHA == "" {
		// 如果没有指定 commit，使用默认分支
		branchName := r.Branch
		if branchName == "" {
			branchName = "main" // 默认分支
		}

		branch, err := client.GetBranch(ctx, r.Owner, r.Name, branchName)
		if err != nil {
			log.Debug().Err(err).Msgf("GitCode: Failed to get branch %s for %s/%s", branchName, r.Owner, r.Name)
			return []*forge_types.FileMeta{}, nil
		}
		commitSHA = branch.Commit.ID
	}

	// 使用 commit SHA 获取目录树（递归获取）
	tree, err := client.GetTree(ctx, r.Owner, r.Name, commitSHA, true)
	if err != nil {
		log.Debug().Err(err).Msgf("GitCode: Failed to get tree for %s/%s at %s", r.Owner, r.Name, commitSHA)
		// 如果获取树失败，尝试直接获取文件内容
		// 对于手动触发的流水线，我们可以返回空的文件列表，让 Woodpecker 使用默认配置
		log.Debug().Msgf("GitCode: Returning empty file list for manual pipeline trigger")
		return []*forge_types.FileMeta{}, nil
	}

	var files []*forge_types.FileMeta

	// 标准化目录路径
	targetDir := strings.TrimPrefix(f, "/")
	if targetDir != "" && !strings.HasSuffix(targetDir, "/") {
		targetDir += "/"
	}

	for _, entry := range tree.Tree {
		// 只处理文件类型
		if entry.Type != "blob" {
			continue
		}

		// 检查文件是否在目标目录中
		if targetDir == "" {
			// 根目录：只要文件不在子目录中
			if !strings.Contains(entry.Path, "/") {
				// 获取文件内容
				data, err := client.GetFileContent(ctx, r.Owner, r.Name, entry.Path, commitSHA)
				if err != nil {
					log.Debug().Err(err).Msgf("GitCode: Failed to get file content for %s", entry.Path)
					continue
				}

				files = append(files, &forge_types.FileMeta{
					Name: entry.Path,
					Data: data,
				})
			}
		} else {
			// 指定目录：文件路径必须以目标目录开头，且不在更深的子目录中
			if strings.HasPrefix(entry.Path, targetDir) {
				relativePath := strings.TrimPrefix(entry.Path, targetDir)
				// 确保文件直接在目标目录中，不在子目录中
				if !strings.Contains(relativePath, "/") {
					// 获取文件内容
					data, err := client.GetFileContent(ctx, r.Owner, r.Name, entry.Path, commitSHA)
					if err != nil {
						log.Debug().Err(err).Msgf("GitCode: Failed to get file content for %s", entry.Path)
						continue
					}

					files = append(files, &forge_types.FileMeta{
						Name: entry.Path,
						Data: data,
					})
				}
			}
		}
	}

	// log.Debug().Msgf("GitCode: Found %d files in directory %s for repo %s", len(files), f, r.FullName)
	return files, nil
}

func (c *GitCode) Status(ctx context.Context, user *model.User, repo *model.Repo, pipeline *model.Pipeline, workflow *model.Workflow) error {
	// GitCode 暂时不支持状态更新，跳过
	// TODO: 实现 GitCode 状态 API 支持
	// log.Debug().Msgf("GitCode status update not implemented for repo %s", repo.FullName)
	return nil
}

func (c *GitCode) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	login := ""
	token := ""

	if u != nil {
		login = u.Login
		token = u.AccessToken
	}

	host, err := common.ExtractHostFromCloneURL(r.Clone)
	if err != nil {
		return nil, err
	}

	return &model.Netrc{
		Login:    login,
		Password: token,
		Machine:  host,
		Type:     model.ForgeTypeGitCode,
	}, nil
}

func (c *GitCode) Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	client := NewGitCodeClient(u.AccessToken, false)

	hook := &CreateHookRequest{
		URL:         link,
		ContentType: "json",
		Events:      []string{"push", "pull_request", "release"},
		Active:      true,
	}

	_, err := client.CreateHook(ctx, r.Owner, r.Name, hook)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("could not find repository")
		}
		return err
	}
	return nil
}

func (c *GitCode) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	client := NewGitCodeClient(u.AccessToken, false)

	hooks, err := client.GetHooks(ctx, r.Owner, r.Name)
	if err != nil {
		return err
	}

	for _, hook := range hooks {
		if hook.URL == link {
			return client.DeleteHook(ctx, r.Owner, r.Name, hook.ID)
		}
	}

	return nil
}

func (c *GitCode) Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	token := common.UserToken(ctx, r, u)
	client := NewGitCodeClient(token, false)

	branches, err := client.GetBranches(ctx, r.Owner, r.Name)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(branches))
	for i := range branches {
		result[i] = branches[i].Name
	}
	return result, err
}

func (c *GitCode) BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (*model.Commit, error) {
	token := common.UserToken(ctx, r, u)
	client := NewGitCodeClient(token, false)

	b, err := client.GetBranch(ctx, r.Owner, r.Name, branch)
	if err != nil {
		return nil, err
	}
	return &model.Commit{
		SHA:      b.Commit.ID,
		ForgeURL: fmt.Sprintf("%s/%s/%s/commit/%s", defaultURL, r.Owner, r.Name, b.Commit.ID),
	}, nil
}

func (c *GitCode) PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error) {
	token := common.UserToken(ctx, r, u)
	client := NewGitCodeClient(token, false)

	pullRequests, err := client.GetPullRequests(ctx, r.Owner, r.Name)
	if err != nil {
		// Repositories without commits return empty list with status code 404
		if strings.Contains(err.Error(), "404") {
			return []*model.PullRequest{}, nil
		}
		return nil, err
	}

	result := make([]*model.PullRequest, len(pullRequests))
	for i := range pullRequests {
		result[i] = convertPullRequest(pullRequests[i])
	}
	return result, err
}

func (c *GitCode) Hook(ctx context.Context, r *http.Request) (*model.Repo, *model.Pipeline, error) {
	repo, pipeline, err := parseHook(r)
	if err != nil {
		return nil, nil, err
	}

	if pipeline != nil && pipeline.Event == model.EventRelease && pipeline.Commit == "" {
		tagName := strings.Split(pipeline.Ref, "/")[2]
		sha, err := c.getTagCommitSHA(ctx, repo, tagName)
		if err != nil {
			return nil, nil, err
		}
		pipeline.Commit = sha
	}

	if pipeline != nil && (pipeline.Event == model.EventPull || pipeline.Event == model.EventPullClosed) && len(pipeline.ChangedFiles) == 0 {
		index, err := strconv.ParseInt(strings.Split(pipeline.Ref, "/")[2], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		pipeline.ChangedFiles, err = c.getChangedFilesForPR(ctx, repo, index)
		if err != nil {
			log.Error().Err(err).Msgf("could not get changed files for PR %s#%d", repo.FullName, index)
		}
	}

	return repo, pipeline, nil
}

func (c *GitCode) OrgMembership(ctx context.Context, u *model.User, owner string) (*model.OrgPerm, error) {
	// GitCode 暂时不支持组织成员检查，返回默认权限
	// TODO: 实现 GitCode 组织成员 API 支持
	return &model.OrgPerm{}, nil
}

func (c *GitCode) Org(ctx context.Context, u *model.User, owner string) (*model.Org, error) {
	// GitCode 暂时不支持组织信息，返回用户信息
	// TODO: 实现 GitCode 组织 API 支持
	return &model.Org{
		Name:    owner,
		IsUser:  true,
		Private: false,
	}, nil
}

// newGitCodeClient 创建新的 GitCode 客户端
func (c *GitCode) newGitCodeClient(token string) *GitCodeClient {
	return NewGitCodeClient(token, false)
}

func (c *GitCode) getChangedFilesForPR(ctx context.Context, repo *model.Repo, index int64) ([]string, error) {
	// GitCode 暂时不支持 PR 文件变更列表，返回空列表
	// TODO: 实现 GitCode PR 文件变更 API 支持
	log.Debug().Msgf("GitCode PR changed files not implemented for repo %s PR #%d", repo.FullName, index)
	return []string{}, nil
}

func (c *GitCode) getTagCommitSHA(ctx context.Context, repo *model.Repo, tagName string) (string, error) {
	// GitCode 暂时不支持标签 API，返回空字符串
	// TODO: 实现 GitCode 标签 API 支持
	log.Debug().Msgf("GitCode tag commit SHA not implemented for repo %s tag %s", repo.FullName, tagName)
	return "", nil
}

func (c *GitCode) perPage(ctx context.Context) int {
	if c.pageSize == 0 {
		c.pageSize = defaultPageSize
	}
	return c.pageSize
}
