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

### GET /api/v1/auth/quota

IMAP storage quota for the current user.

**Response 200**

```json
{"used": 1024, "limit": 1048576}
```

---

## Implemented (P1 — mail)

All endpoints below require authentication (`gsn_session`).

### Folders

| Method | Path | Description |
|--------|------|-------------|
| GET | `/folders` | Folder tree with unread counts |
| POST | `/folders` | Create subfolder (`name`, optional `parent`, `delim`) |
| POST | `/folders/rename` | Rename folder (`name`, `newname`) |
| POST | `/folders/delete` | Delete folder (`name`; blocks system folders) |
| GET | `/folders/:name/count` | Unread count for one folder |

### Messages

| Method | Path | Description |
|--------|------|-------------|
| GET | `/mail/:mailbox` | Paginated list (`page`, `limit`, `q`, `unseen`, `flagged`) |
| GET | `/mail/:mailbox/:uid` | Full message (sanitized HTML, plain, attachments) |
| GET | `/mail/:mailbox/:uid/download` | Download as `.eml` |
| GET | `/mail/:mailbox/:uid/raw` | Raw RFC822 source |
| GET | `/mail/:mailbox/:uid/attachment/:part` | Download attachment by MIME part |
| POST | `/mail/:mailbox/:uid/flag` | Set flag (`seen`, `flagged`, `answered`; `value=1\|0`) |
| POST | `/mail/:mailbox/:uid/move` | Move to folder (`dest`) |
| DELETE | `/mail/:mailbox/:uid` | Move to trash (or expunge if already in trash) |
| DELETE | `/mail/:mailbox` | Empty trash or move all messages to trash |

### Compose

| Method | Path | Description |
|--------|------|-------------|
| POST | `/compose/send` | Send email (`to`, `subject`, `body_html`, attachments) |
| POST | `/compose/draft` | Save draft to IMAP Drafts folder |
| POST | `/compose/upload` | Upload temp attachment (`file` field) |

### Search

| Method | Path | Description |
|--------|------|-------------|
| GET | `/search` | Search mailbox (`q`, optional `mailbox`, `unseen`) |

---

## Planned (P2 frontend)

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
