# StudyHub

A full-stack web application for managing academic modules, sharing study resources, and generating AI-powered flashcards. Students can upload files and links organized by modules, weeks, and academic terms -- and the system automatically generates flashcards from uploaded documents using Google Gemini AI.

## Features

- **Module Management** -- Create and organize academic modules with semester-based runs and weekly structure
- **Resource Sharing** -- Upload files (stored in AWS S3 with deduplication) or share links, organized by week
- **AI Flashcard Generation** -- Uploaded documents are automatically processed by Google Gemini to generate study flashcards
- **Interactive Study Mode** -- Flip-card UI with keyboard navigation for reviewing generated flashcards
- **User Profiles** -- View resources uploaded by any user with full module context
- **Academic Terms** -- Manage semesters and track the active term
- **Admin Dashboard** -- Overview stats, module/run management for administrators
- **Authentication** -- JWT-based auth with role support (admin/regular user)

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | React 18, TypeScript, Vite, Tailwind CSS, shadcn/ui |
| **Backend** | Go, Chi router, PostgreSQL, pgx |
| **Storage** | AWS S3 (file storage with presigned URLs) |
| **AI** | Google Gemini 2.5 Flash (flashcard generation) |
| **Queue** | RabbitMQ (async document processing) |
| **Doc Conversion** | Gotenberg (file-to-PDF conversion) |
| **Containerization** | Docker, Docker Compose |
| **CI/CD** | GitHub Actions, GHCR, EC2 deployment |

## Architecture

```
                    +-----------+
                    |  Frontend |  React SPA (Nginx)
                    +-----+-----+
                          |
                     /api/v1/*
                          |
                    +-----+-----+
                    |  Backend  |  Go API (Chi)
                    +-----+-----+
                       /  |  \
                      /   |   \
              +------+ +--+--+ +----------+
              | PostgreSQL   | |  AWS S3   |
              +--------------+ +----------+
                      |
                 +----+----+
                 | RabbitMQ |
                 +----+----+
                      |
               +------+-------+
               | Worker Pool  |  (5 goroutines)
               +------+-------+
                  /         \
          +------+---+ +----+------+
          | Gotenberg | |  Gemini   |
          | (PDF)     | |  (AI)     |
          +----------+ +-----------+
```

**Flow:** File upload -> S3 storage -> RabbitMQ message -> Worker converts to PDF (via Gotenberg if needed) -> Gemini generates flashcards -> Stored in DB.

## Project Structure

```
StudyHub/
├── backend/
│   ├── cmd/main.go                  # Entry point
│   └── internal/
│       ├── http/                    # Handlers, routing, middleware
│       ├── auth/                    # JWT auth, login, bcrypt
│       ├── modules/                 # Modules, runs, weeks, terms
│       ├── users/                   # User management
│       ├── resources/               # File/link resources, dedup
│       ├── content/                 # Flashcard generation workers
│       ├── aws/                     # S3 storage abstraction
│       ├── gemini/                  # Gemini AI client
│       ├── rabbitmq/                # RabbitMQ client
│       └── config/                  # Environment config
├── frontend/
│   └── src/
│       ├── pages/                   # Page components (10 pages)
│       ├── components/              # UI components (layout, resources, admin, etc.)
│       ├── api/                     # API client layer (Axios)
│       ├── context/                 # Auth context (React Context)
│       └── types/                   # TypeScript type definitions
├── migrations/                      # PostgreSQL migration files
├── docs/                            # API documentation
├── compose.yaml                     # Docker Compose (production)
└── compose.dev.yaml                 # Docker Compose (development)
```

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- AWS account with an S3 bucket
- [Google Gemini API key](https://ai.google.dev/)

### Environment Variables

Create a `.env` file in the project root:

```env
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=studyhub

# JWT
JWT_KEY=your-jwt-secret-key

# AWS S3
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET=your-bucket-name
AWS_DEFAULT_REGION=us-east-1
AWS_S3_URL=https://your-bucket-name.s3.us-east-1.amazonaws.com

# RabbitMQ
RBMQ_USER=guest
RBMQ_PASS=guest
RBMQ_HOST=rabbitmq

# Google Gemini AI
GEMINI_API_KEY=your-gemini-api-key
```

### Run with Docker (Production)

```bash
docker compose up --build
```

This starts all 6 services: frontend (port 80), backend (port 8080), PostgreSQL, RabbitMQ, Gotenberg, and runs database migrations automatically.

### Run with Docker (Development)

```bash
docker compose -f compose.yaml -f compose.dev.yaml up --build
```

Development mode mounts source code as volumes for hot-reloading.

### Run Locally (without Docker)

**Backend:**

```bash
cd backend
go run cmd/main.go
```

Requires Go 1.25+, a running PostgreSQL instance, RabbitMQ, and environment variables configured.

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

The dev server runs on `http://localhost:5173` by default.

## API Overview

All endpoints are under `/api/v1`. Public routes:

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/auth/login` | Authenticate and receive JWT |
| `POST` | `/users` | Register a new user |
| `GET` | `/health` | Health check |

Protected routes (require `Authorization: Bearer <token>`):

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/modules` | List all modules |
| `POST` | `/modules` | Create a module |
| `GET` | `/modules/{id}` | Get module with runs and weeks |
| `POST` | `/resources/file/{week_id}` | Upload a file resource |
| `POST` | `/resources/link/{week_id}` | Create a link resource |
| `GET` | `/resources/weeks/{week_id}` | List resources for a week |
| `GET` | `/resources/{id}` | Get presigned download URL |
| `GET` | `/resources/users/{user_id}` | List resources by user |
| `GET` | `/users/me` | Get current user |

See [`docs/api.v1.md`](docs/api.v1.md) for full API documentation.

## Database Schema

8 tables managed via [golang-migrate](https://github.com/golang-migrate/migrate):

- **modules** -- Academic modules (code, name, department)
- **module_runs** -- Semester instances of modules
- **weeks** -- Weekly structure within runs
- **academic_terms** -- Semester/year tracking
- **users** -- User accounts with bcrypt passwords
- **storage_objects** -- S3 objects with SHA256 hash deduplication
- **resources** -- Files and links with ownership tracking
- **flashcards** -- AI-generated question/answer pairs

## Deployment

The project includes a CI/CD pipeline that:

1. Detects changes in `backend/`, `frontend/`, or `migrations/`
2. Builds and pushes Docker images to GitHub Container Registry
3. Deploys to an EC2 instance via SSH

## License

This project is not currently licensed.
