# ecommerce-api

A REST API built with Go, Gin, GORM, and PostgreSQL.

## Requirements

- Go 1.23+
- PostgreSQL

## Setup

```bash
cp .env.example .env
go mod download
go run ./cmd/server
```

### Environment variables

| Variable | Description |
|----------|-------------|
| `PORT` | Port to listen on (default: `8080`) |
| `DATABASE_URL` | Postgres connection string |
| `JWT_SECRET` | Secret key for signing JWTs |

## Endpoints

### Auth

```
POST /api/v1/register   { "email": "", "password": "" }
POST /api/v1/login      { "email": "", "password": "" }
```

### Products (JWT required)

```
GET    /api/v1/products
GET    /api/v1/products/:id
POST   /api/v1/products      { "name": "", "description": "", "price": 0.0, "stock": 0 }
PUT    /api/v1/products/:id  { "name": "", "description": "", "price": 0.0, "stock": 0 }
DELETE /api/v1/products/:id
```

Pass the token from login as a Bearer token:

```
Authorization: Bearer <token>
```
