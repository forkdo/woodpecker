package gitcode

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePullRequestHook(t *testing.T) {
	// 使用你提供的实际 GitCode webhook 数据
	webhookData := `{
		"enterprise_labels": [],
		"changes": {
			"merge_params": {"current": "force_remove_source_branch: false"},
			"patchset_locked": {"current": false},
			"merge_when_pipeline_succeeds": {"current": false},
			"iid": {"current": 4},
			"target_branch": {"current": "main"},
			"created_at": {"current": "2025-10-01T17:52:19+08:00"},
			"description": {"current": "test"},
			"close_issue_when_merge": {"current": false},
			"moderation_result": {"current": false},
			"source_project_id": {"current": 7720285},
			"title": {"current": "test"},
			"current_patchset_id": {"current": 0},
			"source_branch": {"current": "dev"},
			"squash": {"current": false},
			"updated_at": {"current": "2025-10-01T17:52:20+08:00"},
			"merge_status": {"current": "unchecked"},
			"moderation_time": {"current": 0},
			"latest_merge_request_diff_id": {"current": 3776995},
			"id": {"current": 7326072},
			"state": {"current": "opened"},
			"author_id": {"current": 143790},
			"target_project_id": {"current": 7720285}
		},
		"project": {
			"path_with_namespace": "jetsung/testci",
			"ssh_url": "git@gitcode.com:jetsung/testci.git",
			"description": "",
			"git_http_url": "https://gitcode.com/jetsung/testci.git",
			"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
			"url": "git@gitcode.com:jetsung/testci.git",
			"http_url": "https://gitcode.com/jetsung/testci.git",
			"web_url": "https://gitcode.com/jetsung/testci",
			"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
			"name": "testci",
			"namespace": "jetsung",
			"visibility_level": 20,
			"default_branch": "main",
			"id": 7720285,
			"homepage": "https://gitcode.com/jetsung/testci"
		},
		"git_commit_no": "",
		"virtual_merge_build": false,
		"merge_request": {
			"merge_when_pipeline_succeeds": false,
			"source": {
				"path_with_namespace": "jetsung/testci",
				"ssh_url": "git@gitcode.com:jetsung/testci.git",
				"description": "",
				"git_http_url": "https://gitcode.com/jetsung/testci.git",
				"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
				"url": "git@gitcode.com:jetsung/testci.git",
				"http_url": "https://gitcode.com/jetsung/testci.git",
				"web_url": "https://gitcode.com/jetsung/testci",
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "testci",
				"namespace": "jetsung",
				"visibility_level": 20,
				"default_branch": "main",
				"id": 7720285,
				"homepage": "https://gitcode.com/jetsung/testci"
			},
			"act": "open",
			"oldrev": "",
			"action": "open",
			"id": 7326072,
			"state": "opened",
			"work_in_progress": false,
			"author": {
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "jetsung",
				"id": 143790,
				"email": "i@jetsung.com",
				"username": "jetsung"
			},
			"update_reason": "",
			"act_desc": "open",
			"target_branch": "main",
			"need_approve": false,
			"need_test": false,
			"tester_list": [],
			"total_time_spent": 0,
			"assignee_list": [],
			"approver_list": [],
			"author_id": 143790,
			"target_project_id": 7720285,
			"conflict": false,
			"last_commit": {
				"author": {
					"name": "Jetsung Chan",
					"email": "jetsungchan@gmail.com"
				},
				"id": "e0f538eaf7ded5a29cac7068497f455300b3a5ae",
				"message": "test",
				"url": "https://gitcode.com/jetsung/testci/commits/detail/e0f538eaf7ded5a29cac7068497f455300b3a5ae",
				"timestamp": "2025-10-01T09:50:05Z"
			},
			"iid": 4,
			"created_at": "2025-10-01T17:52:19+08:00",
			"description": "test",
			"title": "test",
			"source_branch": "dev",
			"target_branch_commit": {
				"author": {
					"name": "jetsung",
					"email": "i@jetsung.com"
				},
				"id": "6a6d29ba8df340a5df8c18bb08ab4ee6626476fd",
				"message": "merge dev into maintestCreated-by: jetsungCommit-by: Jetsung ChanMerged-by: jetsungDescription: testSee merge request: jetsung/testci!3",
				"url": "https://gitcode.com/jetsung/testci/commits/detail/6a6d29ba8df340a5df8c18bb08ab4ee6626476fd",
				"timestamp": "2025-10-01T09:33:54Z"
			},
			"need_review": false,
			"updated_at": "2025-10-01T17:52:20+08:00",
			"merge_params": {
				"force_remove_source_branch": false
			},
			"source_project_id": 7720285,
			"url": "https://gitcode.com/jetsung/testci/merge_requests/4",
			"target": {
				"path_with_namespace": "jetsung/testci",
				"ssh_url": "git@gitcode.com:jetsung/testci.git",
				"description": "",
				"git_http_url": "https://gitcode.com/jetsung/testci.git",
				"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
				"url": "git@gitcode.com:jetsung/testci.git",
				"http_url": "https://gitcode.com/jetsung/testci.git",
				"web_url": "https://gitcode.com/jetsung/testci",
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "testci",
				"namespace": "jetsung",
				"visibility_level": 20,
				"default_branch": "main",
				"id": 7720285,
				"homepage": "https://gitcode.com/jetsung/testci"
			},
			"merge_status": "unchecked",
			"reviewer_list": []
		},
		"git_branch": "",
		"repository": {
			"name": "testci",
			"description": "",
			"visibility_level": 20,
			"git_http_url": "https://gitcode.com/jetsung/testci.git",
			"url": "git@gitcode.com:jetsung/testci.git",
			"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
			"homepage": "https://gitcode.com/jetsung/testci"
		},
		"issues": [],
		"object_kind": "merge_request",
		"labels": [],
		"produce_random_id": "47b3ac73567946fa8f63942e36308c8a",
		"event_type": "merge_request",
		"object_attributes": {
			"merge_when_pipeline_succeeds": false,
			"source": {
				"path_with_namespace": "jetsung/testci",
				"ssh_url": "git@gitcode.com:jetsung/testci.git",
				"description": "",
				"git_http_url": "https://gitcode.com/jetsung/testci.git",
				"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
				"url": "git@gitcode.com:jetsung/testci.git",
				"http_url": "https://gitcode.com/jetsung/testci.git",
				"web_url": "https://gitcode.com/jetsung/testci",
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "testci",
				"namespace": "jetsung",
				"visibility_level": 20,
				"default_branch": "main",
				"id": 7720285,
				"homepage": "https://gitcode.com/jetsung/testci"
			},
			"act": "open",
			"oldrev": "",
			"action": "open",
			"id": 7326072,
			"state": "opened",
			"work_in_progress": false,
			"author": {
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "jetsung",
				"id": 143790,
				"email": "i@jetsung.com",
				"username": "jetsung"
			},
			"update_reason": "",
			"act_desc": "open",
			"target_branch": "main",
			"need_approve": false,
			"need_test": false,
			"tester_list": [],
			"total_time_spent": 0,
			"assignee_list": [],
			"approver_list": [],
			"author_id": 143790,
			"target_project_id": 7720285,
			"conflict": false,
			"last_commit": {
				"author": {
					"name": "Jetsung Chan",
					"email": "jetsungchan@gmail.com"
				},
				"id": "e0f538eaf7ded5a29cac7068497f455300b3a5ae",
				"message": "test",
				"url": "https://gitcode.com/jetsung/testci/commits/detail/e0f538eaf7ded5a29cac7068497f455300b3a5ae",
				"timestamp": "2025-10-01T09:50:05Z"
			},
			"iid": 4,
			"created_at": "2025-10-01T17:52:19+08:00",
			"description": "test",
			"title": "test",
			"source_branch": "dev",
			"target_branch_commit": {
				"author": {
					"name": "jetsung",
					"email": "i@jetsung.com"
				},
				"id": "6a6d29ba8df340a5df8c18bb08ab4ee6626476fd",
				"message": "merge dev into maintestCreated-by: jetsungCommit-by: Jetsung ChanMerged-by: jetsungDescription: testSee merge request: jetsung/testci!3",
				"url": "https://gitcode.com/jetsung/testci/commits/detail/6a6d29ba8df340a5df8c18bb08ab4ee6626476fd",
				"timestamp": "2025-10-01T09:33:54Z"
			},
			"need_review": false,
			"updated_at": "2025-10-01T17:52:20+08:00",
			"merge_params": {
				"force_remove_source_branch": false
			},
			"source_project_id": 7720285,
			"url": "https://gitcode.com/jetsung/testci/merge_requests/4",
			"target": {
				"path_with_namespace": "jetsung/testci",
				"ssh_url": "git@gitcode.com:jetsung/testci.git",
				"description": "",
				"git_http_url": "https://gitcode.com/jetsung/testci.git",
				"git_ssh_url": "git@gitcode.com:jetsung/testci.git",
				"url": "git@gitcode.com:jetsung/testci.git",
				"http_url": "https://gitcode.com/jetsung/testci.git",
				"web_url": "https://gitcode.com/jetsung/testci",
				"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
				"name": "testci",
				"namespace": "jetsung",
				"visibility_level": 20,
				"default_branch": "main",
				"id": 7720285,
				"homepage": "https://gitcode.com/jetsung/testci"
			},
			"merge_status": "unchecked",
			"reviewer_list": []
		},
		"git_target_branch_commit_no": "6a6d29ba8df340a5df8c18bb08ab4ee6626476fd",
		"user": {
			"avatar_url": "https://cdn-img.gitcode.com/bf/ee/b4f489e3933733e085b0f6ad0073345142d939d9e53c67d7e9445c66b71ad0a0.JPG?time=1705574695777",
			"name": "jetsung",
			"id": 143790,
			"email": "i@jetsung.com",
			"username": "jetsung"
		},
		"manual_build": false,
		"uuid": "4_d3e9792e-1d6d-499b-ad08-96d60d739469"
	}`

	reader := strings.NewReader(webhookData)
	repo, pipeline, err := parsePullRequestHook(reader)

	// 测试解析是否成功
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NotNil(t, pipeline)

	// 测试 repo 信息
	assert.Equal(t, "7720285", string(repo.ForgeRemoteID))
	assert.Equal(t, "jetsung", repo.Owner)
	assert.Equal(t, "testci", repo.Name)
	assert.Equal(t, "jetsung/testci", repo.FullName)
	assert.Equal(t, "https://gitcode.com/jetsung/testci", repo.ForgeURL)
	assert.Equal(t, "https://gitcode.com/jetsung/testci.git", repo.Clone)
	assert.Equal(t, "git@gitcode.com:jetsung/testci.git", repo.CloneSSH)
	assert.Equal(t, "main", repo.Branch)

	// 测试 pipeline 信息
	assert.Equal(t, "test", pipeline.Title)
	assert.Equal(t, "test", pipeline.Message)
	assert.Equal(t, "jetsung", pipeline.Author)
	assert.Equal(t, "jetsung", pipeline.Sender)
	assert.Equal(t, "i@jetsung.com", pipeline.Email)
	assert.Equal(t, "e0f538eaf7ded5a29cac7068497f455300b3a5ae", pipeline.Commit)
	assert.Equal(t, "refs/pull/4/head", pipeline.Ref)
	assert.Equal(t, "main", pipeline.Branch)
	assert.Equal(t, "dev:main", pipeline.Refspec)
	assert.Equal(t, "https://gitcode.com/jetsung/testci/merge_requests/4", pipeline.ForgeURL)
}
