# REST API

Base path: **`/api/v1`**

All mutating requests (`POST`, `PUT`, `PATCH`, `DELETE`) require CSRF protection:

1. `GET /` (or any page) sets cookie `csrf_token`
2. Send header `X-CSRF-Token` with the same value on mutating requests

Session cookie: **`gsn_session`** (HttpOnly).

## Implemented (P0)

### GET /api/v1/version

Public build info.

**Response 200**

```json
{"version": "0.1.0", "app": "go-snappymail"}
```

---

### POST /api/v1/auth/login

Authenticate against IMAP and create a session.

**Rate limit:** 10 requests/minute per IP.

**Content-Type:** `application/x-www-form-urlencoded`

| Field | Required | Description |
|-------|----------|-------------|
| `username` | yes | Email address |
| `password` | yes | IMAP password |
| `imap_host` | no | Override IMAP host (defaults to config) |

**Response 200**

```json
{"username": "user@test.local"}
```

Sets cookie `gsn_session`.

**Errors**

| Status | Body |
|--------|------|
| 400 | `{"error": "Username and password are required."}` |
| 401 | `{"error": "Invalid credentials or server unreachable."}` |
| 429 | Rate limit exceeded |

**Example**

```bash
# 1. Fetch CSRF token
curl -c cookies.txt http://localhost:8082/

# 2. Login (replace TOKEN from csrf_token cookie)
curl -b cookies.txt -c cookies.txt \
  -H "X-CSRF-Token: TOKEN" \
  -d "username=user@test.local&password=Password1@" \
  http://localhost:8082/api/v1/auth/login
```

---

### POST /api/v1/auth/logout

Invalidate session and clear cookie.

**Response 200**

```json
{"ok": true}
```

---

### GET /api/v1/auth/me

Current session (requires `gsn_session`).

**Response 200**

```json
{
  "username": "user@test.local",
  "datetime_format": "02/01/2006 15:04"
}
```

**Errors:** 401 if not authenticated.

---

## Planned (P1+)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/auth/quota` | IMAP storage quota (handler exists, route pending) |
| GET | `/folders` | List IMAP folders |
| GET | `/folders/{name}/messages` | Message list |
| GET | `/messages/{uid}` | Fetch message |
| POST | `/messages/send` | Compose / send |
| GET | `/search` | IMAP search |

See [OpenSpec](../openspec/changes/go-snappymail-foundation/) for full API design.

## Static SPA

| Path | Behavior |
|------|----------|
| `GET /*` | Serves `web/dist/index.html` (SPA fallback) |
| `GET /assets/*` | Embedded static files |

## Error format

JSON object with `error` string field:

```json
{"error": "Human-readable message."}
```
