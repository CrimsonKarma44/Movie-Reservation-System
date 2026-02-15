# 🎬 Movie Reservation System

A comprehensive backend service for a movie reservation platform built with Go, featuring secure user authentication, JWT-based authorization, and complete CRUD operations for movies, theaters, showtimes, and reservations.

> **🗺️ Project Roadmap:** [roadmap.sh](https://roadmap.sh/projects/movie-reservation-system)

---

## 📋 Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Security Features](#security-features)
- [Usage Examples](#usage-examples)
- [Database Schema](#database-schema)
- [Rate Limiting](#rate-limiting)
- [Error Handling](#error-handling)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## ✨ Features

### Core Functionality
- **User Authentication & Authorization**
  - Secure user registration with email validation
  - Login with JWT token generation
  - Token refresh mechanism with refresh tokens
  - User logout with token cleanup
  - Role-based access control (Admin/User)

- **Movie Management**
  - Create, read, update, and delete movies
  - Movie details including title, description, genre, duration, and rating
  - Admin-only operations

- **Theater Management**
  - Create, read, update, and delete theaters
  - Theater information with location and capacity details
  - Admin-only operations

- **Showtime Management**
  - Create, read, update, and delete showtimes
  - Link movies to specific theaters
  - Schedule management with date and time
  - Admin-only operations

- **Reservation Management**
  - Book movie seats/reservations
  - View reservation history
  - Cancel reservations
  - Real-time availability tracking

### Security Features
- **Authentication & Authorization**
  - JWT (JSON Web Tokens) with separate access and refresh tokens
  - Token expiration and renewal mechanisms
  - Secure token storage with thread-safe implementation
  - Algorithm validation (HS256 only)

- **Access Control**
  - Role-based middleware (Admin verification)
  - Protected endpoints requiring authentication
  - Authorization checks on sensitive operations

- **Rate Limiting**
  - Token bucket algorithm implementation
  - Configurable rate limits per endpoint category
  - Auth endpoints: 5 req/sec with burst of 20
  - Admin endpoints: 10 req/sec with burst of 50
  - Reservation endpoints: 15 req/sec with burst of 100

- **Password Security**
  - Strong password validation
  - Minimum 12 characters required
  - Must contain uppercase, lowercase, digits, and symbols
  - Bcrypt hashing for storage

- **Email Validation**
  - RFC 5322 email format validation
  - Prevents invalid email registrations

- **Security Headers**
  - X-Content-Type-Options: nosniff (MIME type protection)
  - X-Frame-Options: DENY (Clickjacking protection)
  - X-XSS-Protection headers

- **CORS Configuration**
  - Configurable Cross-Origin Resource Sharing
  - Origin validation and whitelisting
  - Secure cookie handling

- **Error Handling**
  - Generic error messages prevent information leakage
  - No sensitive data in error responses
  - Timing attack prevention

---

## 🛠 Technology Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.25.5 |
| **Database** | PostgreSQL |
| **ORM** | GORM |
| **Authentication** | JWT (golang-jwt/jwt/v5) |
| **Cryptography** | golang.org/x/crypto |
| **Environment** | godotenv |
| **HTTP Server** | net/http |

**Dependencies:**
```
github.com/golang-jwt/jwt/v5 v5.3.0
github.com/joho/godotenv v1.5.1
golang.org/x/crypto v0.46.0
gorm.io/driver/postgres v1.6.0
gorm.io/gorm v1.31.1
```

---

## 📁 Project Structure

```
Movie-Reservation-System/
├── cmd/
│   └── movie-reservation-system/
│       └── main.go                 # Application entry point
├── server/
│   ├── app.go                      # App initialization and setup
│   └── urls.go                     # Route definitions and server setup
├── handlers/
│   ├── authHandler.go              # Authentication endpoints
│   ├── movieHandler.go             # Movie management endpoints
│   ├── theaterHandler.go           # Theater management endpoints
│   ├── showtimeHandler.go          # Showtime management endpoints
│   └── reservationHandler.go       # Reservation endpoints
├── middleware/
│   ├── jwtMiddleware.go            # JWT authentication middleware
│   ├── adminMiddleware.go          # Admin role verification
│   ├── corsMiddleware.go           # CORS configuration
│   ├── rateLimitMiddleware.go      # Rate limiting implementation
│   └── securityHeadersMiddleware.go # Security headers
├── models/
│   ├── usersModel.go               # User data structure
│   ├── movieModel.go               # Movie data structure
│   ├── theaterModel.go             # Theater data structure
│   ├── showtimeModel.go            # Showtime data structure
│   ├── reservationModel.go         # Reservation data structure
│   ├── jwtModel.go                 # JWT token structures
│   ├── tokenStoreModel.go          # Token storage
│   ├── envModel.go                 # Environment configuration
│   ├── customErrorModel.go         # Custom error types
│   └── urlModel.go                 # URL/Route models
├── services/
│   ├── authService.go              # Authentication business logic
│   ├── databaseService.go          # Database initialization
│   ├── modelService.go             # Generic CRUD service
│   └── userService.go              # User management service
├── utils/
│   ├── util.go                     # Utility functions
│   ├── emailValidator.go           # Email validation
│   ├── passwordValidator.go        # Password validation
│   ├── context.go                  # Context utilities
│   └── securityAudit.go            # Security audit tools
├── api_test/
│   └── Movie Reservation Service/  # Bruno API test collections
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
└── README.md                       # This file
```

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.25.5 or higher
- **PostgreSQL** 12 or higher
- **Git**

Optional but recommended:
- **Docker** (for containerized setup)
- **Bruno** (API testing - included in repo)

---

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/Movie-Reservation-System.git
cd Movie-Reservation-System
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Database Setup

Create a PostgreSQL database:

```bash
createdb movie_reservation_system
```

Or using SQL:
```sql
CREATE DATABASE movie_reservation_system;
```

### 4. Environment Configuration

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Then edit `.env` with your configuration (see Configuration section below).

### 5. Run the Application

```bash
go run cmd/movie-reservation-system/main.go
```

The server will start and display:
```
API Server Initializing...
Server started on :8080
```

---

## ⚙️ Configuration

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=movie_reservation_system
DB_USER=postgres
DB_PASSWORD=your_secure_password

# JWT Configuration
# Must be at least 32 characters long, different for access and refresh tokens
JWT_ACCESS_SECRET=your_super_secret_access_key_at_least_32_chars_long
JWT_REFRESH_SECRET=your_super_secret_refresh_key_at_least_32_chars_long

# Token Expiration
JWT_ACCESS_EXPIRATION=15m      # Access token validity (15 minutes)
JWT_REFRESH_EXPIRATION=7d      # Refresh token validity (7 days)

# Server Configuration
PORT=8080
ENVIRONMENT=development        # development, staging, or production

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Admin User (Created on first run)
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_PASSWORD=SecureAdminPassword123!
```

**Security Notes:**
- JWT secrets must be at least 32 characters long
- Use different secrets for access and refresh tokens
- Never commit `.env` file to version control
- In production, use environment-specific secrets

---

## 🔌 API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | ❌ No |
| POST | `/api/auth/login` | Login user | ❌ No |
| POST | `/api/auth/logout` | Logout user | ✅ Yes |
| POST | `/api/auth/refresh` | Refresh access token | ✅ Yes |

### Movie Endpoints

| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|-----------|
| GET | `/api/movies` | List all movies | ❌ No | ❌ No |
| GET | `/api/movies/:id` | Get movie details | ❌ No | ❌ No |
| POST | `/api/movies` | Create movie | ✅ Yes | ✅ Yes |
| PUT | `/api/movies/:id` | Update movie | ✅ Yes | ✅ Yes |
| DELETE | `/api/movies/:id` | Delete movie | ✅ Yes | ✅ Yes |

### Theater Endpoints

| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|-----------|
| GET | `/api/theaters` | List all theaters | ❌ No | ❌ No |
| GET | `/api/theaters/:id` | Get theater details | ❌ No | ❌ No |
| POST | `/api/theaters` | Create theater | ✅ Yes | ✅ Yes |
| PUT | `/api/theaters/:id` | Update theater | ✅ Yes | ✅ Yes |
| DELETE | `/api/theaters/:id` | Delete theater | ✅ Yes | ✅ Yes |

### Showtime Endpoints

| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|-----------|
| GET | `/api/showtimes` | List all showtimes | ❌ No | ❌ No |
| GET | `/api/showtimes/:id` | Get showtime details | ❌ No | ❌ No |
| POST | `/api/showtimes` | Create showtime | ✅ Yes | ✅ Yes |
| PUT | `/api/showtimes/:id` | Update showtime | ✅ Yes | ✅ Yes |
| DELETE | `/api/showtimes/:id` | Delete showtime | ✅ Yes | ✅ Yes |

### Reservation Endpoints

| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|-----------|
| GET | `/api/reservations` | List user reservations | ✅ Yes | ❌ No |
| GET | `/api/reservations/:id` | Get reservation details | ✅ Yes | ❌ No |
| POST | `/api/reservations` | Create reservation | ✅ Yes | ❌ No |
| DELETE | `/api/reservations/:id` | Cancel reservation | ✅ Yes | ❌ No |

---

## 🔐 Authentication

### JWT Token Structure

The system uses two types of tokens:

1. **Access Token**
   - Short-lived (default: 15 minutes)
   - Used for API requests
   - Included in Authorization header
   - Format: `Bearer <access_token>`

2. **Refresh Token**
   - Long-lived (default: 7 days)
   - Used to obtain new access tokens
   - Stored securely
   - Can be renewed

### Authentication Flow

```
1. User Registration
   POST /api/auth/register
   Body: { email, password }
   Response: { message, success }

2. User Login
   POST /api/auth/login
   Body: { email, password }
   Response: { accessToken, refreshToken, message }

3. API Request with Token
   GET /api/movies
   Header: Authorization: Bearer <access_token>

4. Token Refresh (when access token expires)
   POST /api/auth/refresh
   Header: Authorization: Bearer <refresh_token>
   Response: { accessToken, message }

5. User Logout
   POST /api/auth/logout
   Header: Authorization: Bearer <access_token>
   Response: { message, success }
```

### Token Usage Example

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"SecurePass123!"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"SecurePass123!"}'

# Response:
# {
#   "accessToken": "eyJhbGciOiJIUzI1NiIs...",
#   "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
#   "message": "Login successful"
# }

# Use token in API requests
curl -X GET http://localhost:8080/api/reservations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## 🔒 Security Features

### 1. Rate Limiting

Protects against brute force and DoS attacks:

```
- Auth Endpoints: 5 requests/second (burst: 20)
- Admin Endpoints: 10 requests/second (burst: 50)
- Reservation Endpoints: 15 requests/second (burst: 100)
```

**Response when rate limit exceeded:**
```
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1613097600
```

### 2. Password Security

Requirements:
- Minimum 12 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one digit (0-9)
- At least one special character (!@#$%^&*)

Example valid passwords:
- `MyPassword123!`
- `Secure@Pass2024`
- `Complex#Pwd99`

### 3. Email Validation

- RFC 5322 format compliance
- Prevents invalid email registrations
- Used for user identification and communication

### 4. CORS Protection

- Configurable allowed origins
- Prevents unauthorized cross-origin requests
- Validates origin headers

### 5. Security Headers

Included in all responses:
- `X-Content-Type-Options: nosniff` - MIME type sniffing prevention
- `X-Frame-Options: DENY` - Clickjacking protection
- `X-XSS-Protection: 1; mode=block` - XSS prevention

### 6. Error Handling

- Generic error messages (no sensitive data leakage)
- No stack traces in production
- Consistent error response format
- Timing attack prevention

### 7. Token Management

- Secure token storage (thread-safe)
- Automatic token cleanup
- Separate secrets for access/refresh tokens
- Algorithm validation (HS256 only)
- Expiration enforcement

---

## 💡 Usage Examples

### Create a Movie (Admin Only)

```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception",
    "description": "A mind-bending thriller",
    "genre": "Sci-Fi",
    "duration": 148,
    "rating": 8.8
  }'
```

### Get All Movies

```bash
curl -X GET http://localhost:8080/api/movies
```

### Create a Reservation

```bash
curl -X POST http://localhost:8080/api/reservations \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "showtimeId": 1,
    "numberOfSeats": 2
  }'
```

### Get User Reservations

```bash
curl -X GET http://localhost:8080/api/reservations \
  -H "Authorization: Bearer <access_token>"
```

### Refresh Access Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Authorization: Bearer <refresh_token>"
```

---

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  is_admin BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Movies Table
```sql
CREATE TABLE movies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  genre VARCHAR(100),
  duration INTEGER,
  rating DECIMAL(3,1),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Theaters Table
```sql
CREATE TABLE theaters (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  location VARCHAR(255),
  capacity INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Showtimes Table
```sql
CREATE TABLE showtimes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  movie_id UUID REFERENCES movies(id),
  theater_id UUID REFERENCES theaters(id),
  show_date DATE NOT NULL,
  show_time TIME NOT NULL,
  available_seats INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Reservations Table
```sql
CREATE TABLE reservations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  showtime_id UUID REFERENCES showtimes(id),
  number_of_seats INTEGER NOT NULL,
  reservation_status VARCHAR(50) DEFAULT 'confirmed',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## ⚡ Rate Limiting

The system implements token bucket rate limiting to protect against abuse:

### Configuration

```go
// Auth endpoints: 5 requests/second, burst of 20
authRateLimiter := middleware.NewRateLimiter(5.0, 20)

// Admin endpoints: 10 requests/second, burst of 50
adminRateLimiter := middleware.NewRateLimiter(10.0, 50)

// Reservation endpoints: 15 requests/second, burst of 100
reservationRateLimiter := middleware.NewRateLimiter(15.0, 100)
```

### Behavior

- Allows burst of requests up to the configured limit
- Refills at the specified rate (requests/second)
- Returns 429 (Too Many Requests) when limit exceeded
- Includes remaining quota in response headers

### Common Scenarios

**Normal Usage:**
```
✅ 3 requests/sec on auth endpoint → Allowed
✅ 7 requests/sec on admin endpoint → Allowed
```

**Rate Limited:**
```
❌ 10 requests/sec on auth endpoint (limit: 5) → 429 Too Many Requests
❌ 50 requests/sec on reservation endpoint (limit: 15) → 429 Too Many Requests
```

---

## ⚠️ Error Handling

The system returns standardized error responses:

### Error Response Format

```json
{
  "error": "Generic error message",
  "statusCode": 400,
  "timestamp": "2024-02-15T10:30:00Z"
}
```

### Common Error Codes

| Status | Error | Cause |
|--------|-------|-------|
| 400 | Bad Request | Invalid input or malformed request |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions (not admin) |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Email already exists |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server-side error |

### Error Handling Best Practices

- Never expose sensitive information in errors
- Provide helpful but generic messages to clients
- Log detailed errors server-side for debugging
- Include proper HTTP status codes
- Consistent error response format

---

## 🛠 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./handlers
go test ./middleware
```

### Building

```bash
# Build binary
go build -o movie-reservation-system cmd/movie-reservation-system/main.go

# Build for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o movie-reservation-system cmd/movie-reservation-system/main.go
```

### Code Standards

- Follow Go conventions (gofmt, golint)
- Write meaningful comments
- Use error wrapping with context
- Test edge cases
- Security-first mindset

### Debugging

Enable debug logging by setting environment variable:

```bash
export DEBUG=true
go run cmd/movie-reservation-system/main.go
```

### API Testing with Bruno

The repository includes Bruno API test collections in `api_test/`:

1. Import the collection in Bruno
2. Set the environment to "Local"
3. Run individual requests or the complete collection

---

## 📝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**
   ```bash
   git commit -m 'Add amazing feature'
   ```
4. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Contribution Guidelines

- Write clear commit messages
- Add tests for new features
- Update documentation
- Follow code standards
- Security review is mandatory

---

## 📚 Additional Resources

- [Security Documentation](./README_SECURITY_FIXES.md)
- [Security Audit Report](./SECURITY_AUDIT_COMPLETE.txt)
- [Rate Limiting Guide](./RATE_LIMITING_GUIDE.md)
- [Troubleshooting Guide](./TROUBLESHOOTING.md)

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🤝 Support

For issues, questions, or suggestions:

- **Issues:** GitHub Issues
- **Documentation:** See README files in root directory
- **Security:** Report security issues privately to security@example.com

---

## 🗓️ Roadmap

Follow the project roadmap: [roadmap.sh](https://roadmap.sh/projects/movie-reservation-system)

### Upcoming Features
- User profiles and preferences
- Advanced search and filtering
- Payment integration
- Email notifications
- Mobile app API enhancements
- Analytics dashboard

---

## 👥 Authors

- **Your Name** - Initial implementation and security hardening

---

**Last Updated:** February 15, 2026  
**Status:** ✅ Production Ready  
**Version:** 1.0.0
