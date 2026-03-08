# GonGo-Simple-Auth

This documentation provides an overview of the **Gongo Simple Auth** project, a lightweight authentication microservice built with Go, Gin, and MongoDB.

---

### Project Overview

The project is designed to be a straightforward, ready-to-use authentication boilerplate. It implements user registration, login with JWT issuance, and a protected "me" route to verify session validity.

### Key Features

* **User Management**: Handles registration and login using MongoDB for persistence.
* **Security**: Uses `bcrypt` for password hashing and HMAC-SHA256 for JWT generation.
* **Middleware Protection**: Includes a `Guard` middleware that validates Bearer tokens for restricted routes.
* **Environment Configuration**: Managed via `.env` files using a custom unmarshaler.
* **Development Tools**: Pre-configured with `Air` for live reloading.

### API Routes (Endpoints)

All routes are prefixed with `/api/v1`.

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| **POST** | `/auth/register` | No | Registers a new user. |
| **POST** | `/auth/login` | No | Authenticates user and returns a JWT. |
| **GET** | `/users/me` | **Yes** | Returns the authenticated user's ID. |

### Technical Stack

* **Language**: Go 1.25.5.
* **Web Framework**: Gin Gonic.
* **Database**: MongoDB Driver.
* **Authentication**: JWT-Go (v5).

### How to Run

1. Configure your `.env` file with `GONGO_MONGO_URI`, `GONGO_PORT` and `GONGO_JWT_SECRET`.
2. Install dependencies: `go mod tidy`.
3. Run with Air: `air` (or `go run cmd/api/main.go`).
