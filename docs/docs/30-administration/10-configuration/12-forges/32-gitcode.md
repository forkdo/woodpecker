# GitCode

Woodpecker comes with built-in support for GitCode.com. To enable GitCode you should configure the Woodpecker container using the following environment variables:

```bash
WOODPECKER_GITCODE=true
WOODPECKER_GITCODE_CLIENT=${GITCODE_CLIENT}
WOODPECKER_GITCODE_SECRET=${GITCODE_SECRET}
```

## Registration

Register your application with GitCode to create your client id and secret. It is very important the authorization callback URL matches your http(s) scheme and hostname exactly with `/authorize` as the path.

Please use this authorization callback URL:

```
https://your-woodpecker-host.com/authorize
```

## Configuration

This is a full list of configuration options. Please note that many of these options use default configuration values that should work for the majority of installations.

```bash
WOODPECKER_GITCODE=true
WOODPECKER_GITCODE_CLIENT=your_gitcode_oauth_client_id
WOODPECKER_GITCODE_SECRET=your_gitcode_oauth_client_secret
```

### GITCODE

- Name: `WOODPECKER_GITCODE`
- Default: `false`

Enables the GitCode driver.

### `WOODPECKER_GITCODE_CLIENT`

> Default: empty

Configures the GitCode OAuth client id. This is used to authorize access.

### `WOODPECKER_GITCODE_SECRET`

> Default: empty

Configures the GitCode OAuth client secret. This is used to authorize access.

## GitCode OAuth Setup

1. Login to your GitCode account
2. Navigate to **Settings** → **Applications** → **OAuth2 Applications**
3. Click **Create New Application**
4. Fill in the application details:
   - **Application Name**: `Woodpecker CI`
   - **Redirect URI**: `https://your-woodpecker-domain.com/authorize`
   - **Scopes**: Select `read:user`, `read:repository`, `write:repository_hook`
5. Save the application and note down the **Client ID** and **Client Secret**

**Important Notes:**

- The redirect URI must exactly match your Woodpecker server's URL with `/authorize` path
- GitCode uses Gitea-compatible OAuth2 endpoints:
  - Authorization URL: `https://gitcode.com/login/oauth/authorize`
  - Token URL: `https://gitcode.com/login/oauth/access_token`
- Make sure to select the correct scopes for proper functionality

## Compatibility

GitCode is based on Gitea and uses Gitea-compatible APIs. Woodpecker uses the Gitea SDK to communicate with GitCode, ensuring full compatibility with:

- OAuth2 authentication
- Repository webhooks
- Pull request events
- Push events
- Tag and release events
- Status reporting
- Git Trees API for file discovery

## API Support

GitCode supports the following APIs that Woodpecker uses:

- **User API**: `/api/v5/user` - Get current user information
- **Repository API**: `/api/v5/user/repos` - List user repositories
- **Git Trees API**: `/api/v5/repos/:owner/:repo/git/trees/:sha` - Get repository file tree
- **File Content API**: `/api/v5/repos/:owner/:repo/raw/:path` - Get file content
- **Branch API**: `/api/v5/repos/:owner/:repo/branches` - List and get branch information
- **Webhook API**: `/api/v5/repos/:owner/:repo/hooks` - Manage repository webhooks

## Limitations

- GitCode must support Gitea API compatibility
- Requires GitCode version 1.18+ for full feature support
- OAuth application must have appropriate permission scopes
