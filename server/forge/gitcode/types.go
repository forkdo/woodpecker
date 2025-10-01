package gitcode

// pullRequestHook GitCode pull request webhook 数据结构
type pullRequestHook struct {
	ObjectKind  string `json:"object_kind"`   // "merge_request"
	EventType   string `json:"event_type"`    // "merge_request"
	GitBranch   string `json:"git_branch"`    // 分支名称
	GitCommitNo string `json:"git_commit_no"` // 提交号
	ManualBuild bool   `json:"manual_build"`  // 是否手动构建
	UUID        string `json:"uuid"`          // 唯一标识符

	User struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`

	Project struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		WebURL            string `json:"web_url"`
		AvatarURL         string `json:"avatar_url"`
		GitSSHURL         string `json:"git_ssh_url"`
		GitHTTPURL        string `json:"git_http_url"`
		Namespace         string `json:"namespace"`
		VisibilityLevel   int    `json:"visibility_level"`
		PathWithNamespace string `json:"path_with_namespace"`
		DefaultBranch     string `json:"default_branch"`
		Homepage          string `json:"homepage"`
		URL               string `json:"url"`
		SSHURL            string `json:"ssh_url"`
		HTTPURL           string `json:"http_url"`
	} `json:"project"`

	MergeRequest struct {
		ID                        int    `json:"id"`
		IID                       int    `json:"iid"`
		Title                     string `json:"title"`
		Description               string `json:"description"`
		State                     string `json:"state"`
		CreatedAt                 string `json:"created_at"`
		UpdatedAt                 string `json:"updated_at"`
		TargetBranch              string `json:"target_branch"`
		SourceBranch              string `json:"source_branch"`
		AuthorID                  int    `json:"author_id"`
		TargetProjectID           int    `json:"target_project_id"`
		SourceProjectID           int    `json:"source_project_id"`
		Action                    string `json:"action"`
		Act                       string `json:"act"`
		ActDesc                   string `json:"act_desc"`
		URL                       string `json:"url"`
		MergeStatus               string `json:"merge_status"`
		WorkInProgress            bool   `json:"work_in_progress"`
		MergeWhenPipelineSucceeds bool   `json:"merge_when_pipeline_succeeds"`

		Author struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"author"`

		LastCommit struct {
			ID        string `json:"id"`
			Message   string `json:"message"`
			URL       string `json:"url"`
			Timestamp string `json:"timestamp"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`

		Source struct {
			ID                int    `json:"id"`
			Name              string `json:"name"`
			PathWithNamespace string `json:"path_with_namespace"`
			WebURL            string `json:"web_url"`
			GitHTTPURL        string `json:"git_http_url"`
			GitSSHURL         string `json:"git_ssh_url"`
			Namespace         string `json:"namespace"`
			VisibilityLevel   int    `json:"visibility_level"`
			DefaultBranch     string `json:"default_branch"`
			AvatarURL         string `json:"avatar_url"`
			Homepage          string `json:"homepage"`
			URL               string `json:"url"`
		} `json:"source"`

		Target struct {
			ID                int    `json:"id"`
			Name              string `json:"name"`
			PathWithNamespace string `json:"path_with_namespace"`
			WebURL            string `json:"web_url"`
			GitHTTPURL        string `json:"git_http_url"`
			GitSSHURL         string `json:"git_ssh_url"`
			Namespace         string `json:"namespace"`
			VisibilityLevel   int    `json:"visibility_level"`
			DefaultBranch     string `json:"default_branch"`
			AvatarURL         string `json:"avatar_url"`
			Homepage          string `json:"homepage"`
			URL               string `json:"url"`
		} `json:"target"`
	} `json:"merge_request"`

	Repository struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		VisibilityLevel int    `json:"visibility_level"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		URL             string `json:"url"`
		Homepage        string `json:"homepage"`
	} `json:"repository"`

	Labels  []interface{} `json:"labels"`
	Changes interface{}   `json:"changes"`
}

// pushHook GitCode push webhook 数据结构
type pushHook struct {
	// 事件基本信息
	ObjectKind  string `json:"object_kind"`  // 事件类型，此处为 "push"
	EventName   string `json:"event_name"`   // 事件名称，与 object_kind 相同
	Before      string `json:"before"`       // 推送前的提交 SHA 值
	After       string `json:"after"`        // 推送后的提交 SHA 值
	Ref         string `json:"ref"`          // Git 引用路径，例如 "refs/heads/main"
	CheckoutSha string `json:"checkout_sha"` // 检出的提交 SHA 值
	Message     string `json:"message"`      // 消息内容

	// 用户信息
	UserID       int    `json:"user_id"`       // 触发推送事件的用户 ID
	UserName     string `json:"user_name"`     // 用户的显示名称
	UserUsername string `json:"user_username"` // 用户的用户名
	UserEmail    string `json:"user_email"`    // 用户的邮箱地址
	UserAvatar   string `json:"user_avatar"`   // 用户头像的 URL 地址

	// 项目信息
	ProjectID int `json:"project_id"` // 项目的唯一标识符
	Project   struct {
		ID                int    `json:"id"`                  // 项目 ID
		Name              string `json:"name"`                // 项目名称
		Description       string `json:"description"`         // 项目描述
		WebURL            string `json:"web_url"`             // 项目的 Web 访问地址
		AvatarURL         string `json:"avatar_url"`          // 项目头像的 URL 地址
		GitSSHURL         string `json:"git_ssh_url"`         // 项目的 SSH 克隆地址
		GitHTTPURL        string `json:"git_http_url"`        // 项目的 HTTP 克隆地址
		Namespace         string `json:"namespace"`           // 项目的命名空间
		VisibilityLevel   int    `json:"visibility_level"`    // 项目可见性级别（0: 私有, 1:公开）
		PathWithNamespace string `json:"path_with_namespace"` // 带命名空间的项目路径
		DefaultBranch     string `json:"default_branch"`      // 项目的默认分支
		Homepage          string `json:"homepage"`            // 项目主页 URL
		URL               string `json:"url"`                 // 项目 Git 仓库 URL
		SSHURL            string `json:"ssh_url"`             // 项目 SSH 克隆 URL
		HTTPURL           string `json:"http_url"`            // 项目 HTTP 克隆 URL
	} `json:"project"`

	// 提交信息
	Commits []struct {
		ID        string `json:"id"`        // 提交的 SHA 值
		Message   string `json:"message"`   // 提交信息
		Timestamp string `json:"timestamp"` // 提交时间戳
		URL       string `json:"url"`       // 提交详情页面的 URL
		Author    struct {
			Name  string `json:"name"`  // 提交作者的名称
			Email string `json:"email"` // 提交作者的邮箱
		} `json:"author"`
		Added    []string `json:"added"`    // 本次提交新增的文件列表
		Removed  []string `json:"removed"`  // 本次提交删除的文件列表
		Modified []string `json:"modified"` // 本次提交修改的文件列表
	} `json:"commits"`

	// 统计信息
	TotalCommitsCount int      `json:"total_commits_count"` // 本次推送的提交总数
	PushOptions       []string `json:"push_options"`        // 推送选项数组

	// 仓库信息
	Repository struct {
		Name            string `json:"name"`             // 仓库名称
		URL             string `json:"url"`              // 仓库 Git URL
		Description     string `json:"description"`      // 仓库描述
		Homepage        string `json:"homepage"`         // 仓库主页 URL
		GitHTTPURL      string `json:"git_http_url"`     // 仓库 HTTP 克隆地址
		GitSSHURL       string `json:"git_ssh_url"`      // 仓库 SSH 克隆地址
		VisibilityLevel int    `json:"visibility_level"` // 仓库可见性级别（0: 私有, 1: 公开）
	} `json:"repository"`

	// 分支和提交信息
	GitBranch   string `json:"git_branch"`    // 推送的目标分支名称
	GitCommitNo string `json:"git_commit_no"` // 最新提交的 SHA 值
	ManualBuild bool   `json:"manual_build"`  // 是否手动构建
	UUID        string `json:"uuid"`          // 本次推送事件的唯一标识符
}

// releaseHook GitCode release webhook 数据结构
type releaseHook struct {
	Action  string      `json:"action"`
	Repo    *Repository `json:"repository"`
	Sender  *User       `json:"sender"`
	Release *Release    `json:"release"`
}

// Release GitCode release 信息
type Release struct {
	ID          int64  `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Author      *User  `json:"author"`
}
