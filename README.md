# MOCBOT Archive API

A RESTful API service for tracking Discord music bot listening sessions and track plays from MOCBOT.

## Architecture

```
cmd/api/          # Application entrypoint
internal/
  ├── database/   # Database layer (GORM models & queries)
  ├── handlers/   # HTTP handlers (Gin)
  ├── middleware/ # Authentication middleware
  ├── models/     # Domain models & request/response DTOs
  └── utils/      # Utilities (secret management, etc.)
```

## Development Setup

### 1. Set Up Secrets

Create a `.local-secrets` directory in the project root with the following files:

```bash
mkdir .local-secrets
```

**`.local-secrets/api-key`**

```
your-api-key-here
```

**`.local-secrets/db-password`**

```
your-mysql-password
```

**`.local-secrets/db-root-password`**

```
your-mysql-root-password
```

> **Note**: These files should contain only the raw secret values (no quotes, no newlines at the end).

### 2. Configure Environment Variables

The `.env.local` file is already configured with sensible defaults:

```dotenv
API_KEY_FILE=/secrets/api-key
MYSQL_DATABASE=mocbot-archive
MYSQL_USER=mocbot
MYSQL_PASSWORD_FILE=/secrets/db-password
MYSQL_ROOT_PASSWORD_FILE=/secrets/db-root-password
MYSQL_HOST=db
```

You can modify these if needed, but the defaults work out of the box with Docker Compose.

### 3. Start the Application

Using Docker Compose:

```bash
docker compose up --build
```

The API will be available at `http://localhost:9000`.

### 4. Verify It's Running

```bash
curl http://localhost:9000/api/v1/health
```

Expected response:

```json
{ "status": "ok" }
```

## API Routes

All routes except `/health` require an `Authorization` header with the API key, stored in `X-API-Key`.

### Health Check

| Method | Endpoint         | Description                     |
| ------ | ---------------- | ------------------------------- |
| `GET`  | `/api/v1/health` | Health check (no auth required) |

### Sessions

#### `GET /api/v1/sessions`

Get sessions by guild IDs

- **Query Params:**
  - `guildIds` (required) - Comma-separated list of guild IDs (e.g., `123,456`)
  - `limit` (optional) - Maximum number of sessions to return (default: 50)

#### `POST /api/v1/sessions`

Start a new listening session

- **Body:**
  ```json
  {
    "guildId": 123456789
  }
  ```

#### `GET /api/v1/sessions/:sessionId`

Get session by ID

- **URL Params:**
  - `sessionId` (required) - Session ID

#### `PATCH /api/v1/sessions/:sessionId`

End a session

- **URL Params:**
  - `sessionId` (required) - Session ID
- **Body (optional):**
  ```json
  {
    "endedAt": "2026-01-16T10:30:00Z"
  }
  ```
  _If `endedAt` is omitted, current time is used_

#### `GET /api/v1/sessions/:sessionId/tracks`

Get all tracks in a session

- **URL Params:**
  - `sessionId` (required) - Session ID

#### `POST /api/v1/sessions/:sessionId/tracks`

Start a track (creates track play + listeners atomically)

- **URL Params:**
  - `sessionId` (required) - Session ID
- **Body:**
  ```json
  {
    "source": "youtube",
    "sourceId": "dQw4w9WgXcQ",
    "title": "Never Gonna Give You Up",
    "artist": "Rick Astley",
    "url": "https://youtube.com/watch?v=dQw4w9WgXcQ",
    "durationMs": 213000,
    "queuedByUser": 123456789,
    "listenerIds": [123456789, 987654321]
  }
  ```

### Tracks

#### `GET /api/v1/tracks/:trackPlayId`

Get track play by ID

- **URL Params:**
  - `trackPlayId` (required) - Track play ID

#### `PATCH /api/v1/tracks/:trackPlayId`

End a track

- **URL Params:**
  - `trackPlayId` (required) - Track play ID
- **Body (optional):**
  ```json
  {
    "endedAt": "2026-01-16T10:08:30Z"
  }
  ```
  _If `endedAt` is omitted, current time is used_

#### `GET /api/v1/tracks/:trackPlayId/listeners`

Get all listeners for a track

- **URL Params:**
  - `trackPlayId` (required) - Track play ID

### Guilds

#### `GET /api/v1/guilds/:guildId/tracks/recent`

Get recent tracks played in a guild

- **URL Params:**
  - `guildId` (required) - Guild ID
- **Query Params:**
  - `limit` (optional) - Maximum number of tracks to return (default: 50, max: 500)

### Users

#### `GET /api/v1/users/:userId/tracks/recent`

Get recent tracks a user listened to

- **URL Params:**
  - `userId` (required) - User ID
- **Query Params:**
  - `limit` (optional) - Maximum number of tracks to return (default: 50, max: 500)
