# StudyHub

A full-stack web application for managing academic modules, sharing study resources, and generating AI-powered flashcards. Students can upload files and links organized by modules, weeks, and academic terms -- and the system automatically generates flashcards from uploaded documents using Google Gemini AI. It also includes a RAG-based chatbot assistant that answers questions about the application using local Ollama models.

## Features

- **Module Management** -- Create and organize academic modules with semester-based runs and weekly structure
- **Resource Sharing** -- Upload files (stored in AWS S3 with deduplication) or share links, organized by week
- **AI Flashcard Generation** -- Uploaded documents are automatically processed by Google Gemini to generate study flashcards
- **Interactive Study Mode** -- Flip-card UI with keyboard navigation for reviewing generated flashcards
- **RAG Chatbot** -- Context-aware help-desk assistant powered by local Ollama models (gemma4:e4b + nomic-embed-text-v2-moe), rejects off-topic questions
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
| **AI (flashcards)** | Google Gemini 2.5 Flash |
| **AI (chatbot)** | Ollama — gemma4:e4b (LLM) + nomic-embed-text-v2-moe (embeddings) |
| **RAG** | Python, FastAPI, LangChain, ChromaDB |
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
                       /  |  \  \
                      /   |   \  \________________
              +------+ +--+--+ +----------+       \
              | PostgreSQL   | |  AWS S3   |  +----+-------+
              +--------------+ +----------+  | RAG Service |
                      |                      | Python/FAST |
                 +----+----+                 +----+--------+
                 | RabbitMQ |                     |
                 +----+----+              +-------+-------+
                      |                  |               |
               +------+-------+    +----------+   +--------+
               | Worker Pool  |    | ChromaDB |   | Ollama |
               +------+-------+    +----------+   +--------+
                  /         \
          +------+---+ +----+------+
          | Gotenberg | |  Gemini   |
          | (PDF)     | |  (AI)     |
          +----------+ +-----------+
```

**Flashcard flow:** File upload → S3 storage → RabbitMQ message → Worker converts to PDF (via Gotenberg if needed) → Gemini generates flashcards → Stored in DB.

**Chatbot flow:** User question → Go backend → RAG service → ChromaDB similarity search → Ollama (gemma4:e4b) → grounded answer + sources.

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
├── rag-service/
│   ├── main.py                      # FastAPI app — /health, /chat, /rebuild
│   ├── rag.py                       # RAGChatbot — ingestion, retrieval, generation
│   ├── build_index.py               # Build-time indexing script
│   ├── requirements.txt             # Python dependencies
│   ├── Dockerfile
│   └── data/
│       └── studyhub_docs.md         # Knowledge base
├── migrations/                      # PostgreSQL migration files
├── docs/                            # API documentation
├── RAG_CHATBOT_DOCS.md              # RAG chatbot technical documentation
├── compose.yaml                     # Docker Compose (production)
└── compose.dev.yaml                 # Docker Compose (development)
```

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- AWS account with an S3 bucket
- [Google Gemini API key](https://ai.google.dev/)
- [Ollama](https://ollama.com) installed and running (for the chatbot)

### Chatbot Setup (Ollama models)

The RAG chatbot requires two local models. Pull them once before starting the application:

```bash
ollama pull nomic-embed-text-v2-moe   # embedding model (~957 MB)
ollama pull gemma4:e4b                 # LLM (~3 GB)
```

Ollama must be running on your machine (`ollama serve`) before starting the Docker containers. The RAG service connects to it via `host.docker.internal:11434`.

On first container start the chatbot will automatically index `rag-service/data/studyhub_docs.md` into ChromaDB. Subsequent starts load the index from disk and are fast.

To add new documents to the knowledge base, drop `.pdf`, `.txt`, or `.md` files into `rag-service/data/` and call:

```bash
curl -X POST http://localhost:8001/rebuild
```

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

# RAG Chatbot
RAG_SERVICE_URL=http://rag-service:8001
OLLAMA_BASE_URL=http://host.docker.internal:11434
```

### Run with Docker (Production)

```bash
docker compose up --build
```

This starts all 7 services: frontend (port 80), backend (port 8080), RAG service (port 8001), PostgreSQL, RabbitMQ, Gotenberg, and runs database migrations automatically.

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
