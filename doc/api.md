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

## Upload

Api for uploading large files.

All requests must have a valid `Authorization` header.

Size of each block is 8MB.

Only the last block can be smaller than 8MB.

### Create Task

POST /api/upload/create

#### Request

```json
{
  "file_name": "filename",
  "block_count": 10,
  "total_size": 83886079,
  "md5": "md5"
}
```

- block_count: number of blocks
- total_size: total size of the file in bytes
- md5: md5 hash of the file

#### Response

```json
{
  "message": "success",
  "upload_id": "upload_id"
}
```

### Upload Block

PUT /api/upload/block

#### Request

Query parameters:
- upload_id: upload id

Use binary data as the request body.

#### Response

```json
{
  "message": "success"
}
```

### Finish Task

POST /api/upload/finish

#### Request

Query parameters:
- upload_id: upload id

Empty request body.

#### Response

```json
{
  "message": "success"
}
```

### Get Status

GET /api/upload/status

#### Request

Query parameters:
- upload_id: upload id

#### Response

```json
{
  "message": "success",
  "status": "status",
  "created_at": "created_at",
  "block_count": 10,
  "total_size": 83886079
}
```

- status: a string consisting of 0 and 1, where 0 means the block is not uploaded and 1 means the block is uploaded
- created_at: timestamp of the creation time

### Cancel Task

POST /api/upload/cancel

#### Request

Query parameters:
- upload_id: upload id

Empty request body.

#### Response

```json
{
  "message": "success"
}
```