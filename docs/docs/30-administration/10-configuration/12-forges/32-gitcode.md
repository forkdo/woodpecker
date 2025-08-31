# GitCode

Woodpecker comes with built-in support for GitCode.com. To enable GitCode you should configure the Woodpecker container using the following environment variables:

```bash
WOODPECKER_GITCODE=true
WOODPECKER_GITCODE_CLIENT_ID=${GITCODE_CLIENT_ID}
WOODPECKER_GITCODE_CLIENT_SECRET=${GITCODE_CLIENT_SECRET}
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
WOODPECKER_GITCODE_CLIENT_ID=your_gitcode_oauth_client_id
WOODPECKER_GITCODE_CLIENT_SECRET=your_gitcode_oauth_client_secret
WOODPECKER_GITCODE_SKIP_VERIFY=false
WOODPECKER_GITCODE_OAUTH_HOST=https://gitcode.com
```

### GITCODE

- Name: `WOODPECKER_GITCODE`
- Default: `false`

Enables the GitCode driver.

### `WOODPECKER_GITCODE_CLIENT_ID`

> Default: empty

Configures the GitCode OAuth client id. This is used to authorize access.

### `WOODPECKER_GITCODE_CLIENT_SECRET`

> Default: empty

Configures the GitCode OAuth client secret. This is used to authorize access.

### `WOODPECKER_GITCODE_SKIP_VERIFY`

> Default: `false`

Configure if SSL verification should be skipped.

### `WOODPECKER_GITCODE_OAUTH_HOST`

> Default: same as `WOODPECKER_GITCODE_URL`

Configure the OAuth host if it differs from the GitCode URL.

## GitCode OAuth Setup

1. Login to your GitCode account
2. Navigate to **Settings** → **Applications** → **OAuth2 Applications**
3. Click **Create New Application**
4. Fill in the application details:
   - **Application Name**: `Woodpecker CI`
   - **Redirect URI**: `https://your-woodpecker-domain.com/authorize`
   - **Scopes**: Select `read:user`, `read:repository`, `write:repository_hook`
5. Save the application and note down the **Client ID** and **Client Secret**

## Compatibility

GitCode is based on Gitea and uses Gitea-compatible APIs. Woodpecker uses the Gitea SDK to communicate with GitCode, ensuring full compatibility with:

- OAuth2 authentication
- Repository webhooks
- Pull request events
- Push events
- Tag and release events
- Status reporting

## Limitations

- GitCode must support Gitea API compatibility
- Requires GitCode version 1.18+ for full feature support
- OAuth application must have appropriate permission scopes
