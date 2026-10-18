# Breakthrough Palace Mission - Backend

Go backend for the family records platform using GORM and PostgreSQL.

## Project Structure

```
backend/
├── config/              # Configuration management
├── database/            # Database initialization and migrations
├── models/              # GORM models
├── handlers/            # HTTP request handlers
├── middleware/          # HTTP middleware (auth, CORS)
├── auth/                # JWT and password utilities
├── main.go              # Entry point
├── go.mod              # Go module file
├── .env.example        # Environment variables template
└── README.md           # This file
```

## Setup

### Prerequisites
- Go 1.23+
- PostgreSQL 12+

### Installation

1. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Update .env with your database credentials:**
   ```bash
   DATABASE_URL=postgres://postgres:password@localhost:5432/breakthrough_dev
   JWT_SECRET=your-secret-key-here
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` and automatically run migrations.

## Database

Automatic migrations create these tables:
- `families` - Family organizations
- `members` - Family members with auth
- `meetings` - Scheduled meetings
- `attendance` - Attendance records
- `contributions` - Financial contributions
- `minutes` - Meeting minutes/documents

## API Endpoints

### Auth (No auth required)
- `POST /api/auth/login` - Login with email/password
- `POST /api/auth/logout` - Logout (clears cookie)

### Protected Routes (Auth required)
- `GET /api/auth/me` - Get current member info
- `GET /api/family` - Get family details
- `GET /api/family/stats` - Get family statistics

## Cookie Auth

Login creates an `HttpOnly` cookie that:
- Is automatically sent with requests (dev: localhost, prod: .breakthroughpalacemission.com)
- Cannot be accessed by JavaScript (secure)
- Lasts 7 days

## Environment Variables

| Variable | Dev Default | Prod |
|----------|------------|------|
| ENV | development | production |
| PORT | 8080 | 8080 |
| DATABASE_URL | localhost | prod-db |
| JWT_SECRET | dev-key | strong-random-key |
| FRONTEND_URL | http://localhost:3000 | https://dashboard.breakthroughpalacemission.com |

## Development

### Create new handler:
```go
// handlers/member.go
func GetMembers(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        familyID := middleware.ExtractFamilyID(r)
        // ... logic
    }
}
```

### Register route in main.go:
```go
mux.HandleFunc("GET /api/members", handlers.GetMembers(db))
```

## Notes

- All queries are filtered by `FamilyID` for multi-tenant isolation
- Passwords are hashed with bcrypt
- JWTs expire after 7 days
- CORS is configured for frontend URL only
