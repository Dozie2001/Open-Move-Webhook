# Intercord Backend

Intercord is a blockchain events subscription platform, focused on the Sui network. The platform allows users to create subscriptions to specific events on the Sui network and receive notifications through various channels (webhook, email, Telegram, Discord, etc.).

## Features

- Authentication
  - Login/Register
  - Email Verification
  - Password Reset
- Team/Organization Management
  - Create/Join/Leave/Delete teams
  - Invite members with different roles
- Event Subscriptions
  - Create/Edit/Delete subscriptions for blockchain events
  - Configure subscription properties
- Notification Channels
  - Create/Edit/Delete notification channels (webhook, email, Telegram, Discord)
  - Subscribe/Unsubscribe channels to/from subscriptions
- Notifications
  - View notification history
  - Get notification details

## Technology Stack

- **Language**: Golang
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: Bun
- **Authentication**: JWT
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go (v1.20 or higher) for local development

### Running with Docker

1. Clone the repository
```bash
git clone https://github.com/your-organization/intercord-backend.git
cd intercord-backend
```

2. Start the application
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Local Development

1. Clone the repository
```bash
git clone https://github.com/your-organization/intercord-backend.git
cd intercord-backend
```

2. Install dependencies
```bash
go mod download
```

3. Start PostgreSQL and MailHog with Docker
```bash
docker-compose up -d postgres mailhog
```

4. Run the application
```bash
go run ./cmd/app
```

## API Documentation

### Authentication Endpoints

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login
- `GET /auth/verify-email` - Verify email
- `POST /auth/request-reset-password` - Request password reset
- `POST /auth/reset-password` - Reset password

### Team Endpoints

- `GET /teams` - List user's teams
- `POST /teams` - Create a team
- `GET /teams/:id` - Get team details
- `DELETE /teams/:id` - Delete a team
- `POST /teams/:id/invite` - Invite a user to a team
- `POST /teams/:id/join` - Join a team
- `POST /teams/:id/leave` - Leave a team
- `GET /teams/:team_id/subscriptions` - Get team subscriptions
- `GET /teams/:team_id/channels` - Get team channels

### Subscription Endpoints

- `GET /subscriptions` - List user's subscriptions
- `POST /subscriptions` - Create a subscription
- `GET /subscriptions/:id` - Get subscription details
- `PUT /subscriptions/:id` - Update a subscription
- `DELETE /subscriptions/:id` - Delete a subscription

### Channel Endpoints

- `GET /channels` - List user's channels
- `POST /channels` - Create a channel
- `GET /channels/:id` - Get channel details
- `PUT /channels/:id` - Update a channel
- `DELETE /channels/:id` - Delete a channel
- `POST /channels/subscribe` - Subscribe a channel to a subscription
- `POST /channels/unsubscribe` - Unsubscribe a channel from a subscription

### Notification Endpoints

- `GET /notifications` - List notifications
- `GET /notifications/:id` - Get notification details

## Environment Variables

- `SERVER_PORT` - Port for the HTTP server (default: 8080)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASSWORD` - PostgreSQL password (default: postgres)
- `DB_NAME` - PostgreSQL database name (default: intercord)
- `DB_SSLMODE` - PostgreSQL SSL mode (default: disable)
- `JWT_SECRET` - Secret for JWT tokens (must be set in production)
- `JWT_ACCESS_TOKEN_TTL` - Access token lifetime (default: 15m)
- `JWT_REFRESH_TOKEN_TTL` - Refresh token lifetime (default: 7d)
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `EMAIL_FROM` - Sender email address
- `EMAIL_NAME` - Sender name
- `BASE_URL` - Base URL for the application (used in email links)

## Security Considerations

- JWT tokens are used for authentication
- Passwords are securely hashed with bcrypt
- Email verification is required for new accounts
- Role-based access control for team operations
- All endpoints (except authentication) require valid JWT token