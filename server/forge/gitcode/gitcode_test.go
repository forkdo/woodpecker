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
	"testing"

	"github.com/stretchr/testify/assert"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func TestGitCode(t *testing.T) {
	opts := Opts{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
	}

	forge, err := New(opts)
	assert.NoError(t, err)
	assert.NotNil(t, forge)

	gitcode, ok := forge.(*GitCode)
	assert.True(t, ok)
	assert.Equal(t, "gitcode", gitcode.Name())
	assert.Equal(t, "https://gitcode.com", gitcode.URL())
}

func TestGitCodeNetrc(t *testing.T) {
	opts := Opts{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
	}

	forge, err := New(opts)
	assert.NoError(t, err)

	user := &model.User{
		Login:       "testuser",
		AccessToken: "test-token",
	}

	repo := &model.Repo{
		Clone: "https://gitcode.com/testuser/testrepo.git",
	}

	netrc, err := forge.Netrc(user, repo)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", netrc.Login)
	assert.Equal(t, "test-token", netrc.Password)
	assert.Equal(t, "gitcode.com", netrc.Machine)
	assert.Equal(t, model.ForgeTypeGitCode, netrc.Type)
}
