# API

## Auth

Set the `Authorization` header to `token` for all requests that require authentication.

## User

### Register

POST /api/user/register

#### Request

```json
{
  "username": "username",
  "password": "password"
}
```

#### Response

```json
{
  "token": "token"
}
```

### Login

POST /api/user/login

#### Request

```json
{
  "username": "username",
  "password": "password"
}
```

#### Response

```json
{
  "token": "token"
}
```

### Get user info

GET /api/user/:username

#### Response

```json
{
  "username": "username",
  "nickname": "nickname",
  "avatar": "avatar",
  "bio": "bio",
  "is_admin": false
}
```

- avatar: path to the avatar image, or null if the user has no avatar. e.g. /api/user/avatar/:username

### Update user info

POST /api/user/update

#### Request

```json
{
  "nickname": "nickname",
  "avatar": "avatar",
  "bio": "bio"
}
```

All fields are optional.

A valid `Authorization` header is required.

- avatar: base64 encoded image data.

#### Response

```json
{
  "message": "success"
}
```

### Get user avatar

GET /api/user/avatar/:username

### Change password

POST /api/user/changePassword

#### Request

```json
{
  "old_password": "old_password",
  "new_password": "new_password"
}
```

A valid `Authorization` header is required.

#### Response

```json
{
  "message": "success"
}
```

### Set admin

POST /api/user/setAdmin

#### Request

```json
{
  "username": "username",
  "is_admin": true
}
```

A valid `Authorization` header is required.
The user must be an admin.

#### Response

```json
{
  "message": "success"
}
```

### Get user list

GET /api/user/list

#### Request

Query parameters:
- page: page number, default 1

A valid `Authorization` header is required.
The user must be an admin.

#### Response

```json
{
  "message": "success",
  "users": [
    {
      "username": "username",
      "nickname": "nickname",
      "avatar": "avatar",
      "bio": "bio",
      "is_admin": false
    }
  ],
  "max_page": 1
}
```