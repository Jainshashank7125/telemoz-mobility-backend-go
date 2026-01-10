# TelemoZ Mobility Backend

Go backend API for TelemoZ Mobility platform supporting delivery, taxi, and school bus services with real-time tracking.

## Features

- **Authentication**: JWT-based authentication with refresh tokens
- **Trip Management**: Create, track, and manage trips for customers
- **Job Management**: Driver job acceptance, tracking, and completion
- **School Bus Tracking**: Real-time bus location tracking for parents
- **Notifications**: SMS, voice calls, and push notifications
- **Earnings**: Driver earnings calculation and reporting

## Technology Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Real-time Tracking**: Traccar integration
- **Maps**: Google Maps API
- **Notifications**: Twilio (SMS/Voice), Firebase (Push)

## Project Structure

```
telemoz-backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── models/          # Database models
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Data access layer
│   ├── middleware/      # HTTP middleware
│   ├── utils/           # Utility functions
│   └── dto/             # Data transfer objects
├── pkg/
│   ├── traccar/         # Traccar client
│   ├── maps/            # Google Maps client
│   ├── sms/             # SMS provider
│   └── voice/           # Voice call provider
└── api/                 # Route definitions
```

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Traccar server (optional, for real-time tracking)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd telemoz-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file from `.env.example`:
```bash
cp .env.example .env
```

4. Update `.env` with your configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=telemoz
DB_PASSWORD=password
DB_NAME=telemoz_db
JWT_SECRET=your-secret-key
GOOGLE_MAPS_API_KEY=your-api-key
```

5. Create PostgreSQL database:
```bash
createdb telemoz_db
```

6. Run the application:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

## Docker Setup

1. Build and run with Docker Compose:
```bash
docker-compose up -d
```

2. View logs:
```bash
docker-compose logs -f backend
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout user

### Trips (Customer)
- `POST /api/trips` - Create trip
- `GET /api/trips/active` - Get active trip
- `GET /api/trips/history` - Get trip history
- `GET /api/trips/:id` - Get trip details
- `PUT /api/trips/:id` - Update trip
- `POST /api/trips/:id/cancel` - Cancel trip

### Jobs (Driver)
- `GET /api/jobs/available` - Get available jobs
- `POST /api/jobs/:id/accept` - Accept job
- `POST /api/jobs/:id/reject` - Reject job
- `GET /api/jobs/active` - Get active job
- `GET /api/jobs/history` - Get job history
- `PUT /api/jobs/:id/status` - Update job status

### Children (Parent)
- `GET /api/children` - List children
- `POST /api/children` - Add child
- `GET /api/children/:id` - Get child details
- `PUT /api/children/:id` - Update child
- `DELETE /api/children/:id` - Delete child

### Bus Tracking
- `GET /api/buses/child/:childId` - Get bus for child
- `GET /api/buses/:id/track` - Get bus location

### Notifications
- `GET /api/notifications` - List notifications
- `PUT /api/notifications/:id/read` - Mark as read
- `GET /api/notifications/settings` - Get settings
- `PUT /api/notifications/settings` - Update settings

### Earnings (Driver)
- `GET /api/earnings/summary` - Get earnings summary
- `GET /api/earnings/history` - Get earnings history

### Profile
- `GET /api/profile` - Get profile
- `PUT /api/profile` - Update profile

## Database Migrations

The application uses GORM AutoMigrate to automatically create/update database schema on startup. For production, consider using a migration tool like `golang-migrate`.

## Environment Variables

See `.env.example` for all available configuration options.

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Building
```bash
go build -o server ./cmd/server
```

## License

Proprietary - All rights reserved

