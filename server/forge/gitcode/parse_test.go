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
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func TestParseHook(t *testing.T) {
	tests := []struct {
		name        string
		hookType    string
		payload     string
		expectRepo  bool
		expectEvent model.WebhookEvent
	}{
		{
			name:        "push hook",
			hookType:    hookPush,
			payload:     pushHookPayload,
			expectRepo:  true,
			expectEvent: model.EventPush,
		},
		{
			name:        "unsupported hook",
			hookType:    "issues",
			payload:     "",
			expectRepo:  false,
			expectEvent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/hook", strings.NewReader(tt.payload))
			req.Header.Set(hookEvent, tt.hookType)

			repo, pipeline, err := parseHook(req)

			if tt.expectRepo {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				if pipeline != nil {
					assert.Equal(t, tt.expectEvent, pipeline.Event)
				}
			} else {
				assert.Nil(t, repo)
				assert.Nil(t, pipeline)
			}
		})
	}
}

const pushHookPayload = `{
  "ref": "refs/heads/main",
  "before": "0000000000000000000000000000000000000000",
  "after": "1234567890abcdef1234567890abcdef12345678",
  "compare_url": "https://gitcode.com/user/repo/compare/0000000000000000000000000000000000000000...1234567890abcdef1234567890abcdef12345678",
  "commits": [
    {
      "id": "1234567890abcdef1234567890abcdef12345678",
      "message": "Initial commit",
      "url": "https://gitcode.com/user/repo/commit/1234567890abcdef1234567890abcdef12345678",
      "author": {
        "name": "Test User",
        "email": "test@example.com"
      },
      "added": ["README.md"],
      "removed": [],
      "modified": []
    }
  ],
  "head_commit": {
    "id": "1234567890abcdef1234567890abcdef12345678",
    "message": "Initial commit",
    "url": "https://gitcode.com/user/repo/commit/1234567890abcdef1234567890abcdef12345678",
    "author": {
      "name": "Test User",
      "email": "test@example.com"
    },
    "added": ["README.md"],
    "removed": [],
    "modified": []
  },
  "repository": {
    "id": 123,
    "name": "repo",
    "full_name": "user/repo",
    "html_url": "https://gitcode.com/user/repo",
    "clone_url": "https://gitcode.com/user/repo.git",
    "ssh_url": "git@gitcode.com:user/repo.git",
    "default_branch": "main",
    "private": false,
    "owner": {
      "id": 456,
      "login": "user",
      "avatar_url": "https://gitcode.com/avatars/456"
    },
    "permissions": {
      "admin": true,
      "push": true,
      "pull": true
    }
  },
  "pusher": {
    "id": 456,
    "login": "user",
    "email": "test@example.com",
    "avatar_url": "https://gitcode.com/avatars/456"
  },
  "sender": {
    "id": 456,
    "login": "user",
    "email": "test@example.com",
    "avatar_url": "https://gitcode.com/avatars/456"
  }
}`
