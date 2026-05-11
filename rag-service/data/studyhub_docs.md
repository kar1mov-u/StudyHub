# StudyHub — Complete System Documentation

## What is StudyHub?

StudyHub is a full-stack academic study management platform built for university students and lecturers. It helps students organise their coursework by providing a structured hierarchy of modules, module runs, weeks, and resources. Students can upload study files, share links, generate AI-powered flashcards from uploaded PDFs, comment on resources, and view other students' contributions.

StudyHub is built and maintained by a student engineering team at Central Asian University as part of the Software Engineering and Internship module.

---

## Core Features

### 1. Modules
A Module represents an academic course (e.g. "Introduction to Computer Science", "Data Structures"). Each module has a department name and a unique code. Modules are listed on the Modules page and can be created by admin users.

### 2. Module Runs
A Module Run is a specific delivery of a module in a given semester and year (e.g. "Spring 2026"). Each run contains weekly sessions. A module can have multiple runs across different academic terms.

### 3. Weeks
Each Module Run is divided into individual Weeks (e.g. Week 1, Week 2, ..., Week 12). Each week holds study resources uploaded by students and teaching staff.

### 4. Resources
Resources are study materials attached to a specific week. There are two types:
- **File resources** — PDFs, Word documents, images, etc. Files are stored in AWS S3. StudyHub deduplicates files using SHA-256 hashing so identical files are stored only once.
- **Link resources** — External URLs (e.g. YouTube lectures, Wikipedia articles, documentation pages).

Resources have an owner (the user who uploaded them), a name, a type (file or link), and a creation timestamp.

### 5. AI Flashcard Generation
When a student uploads a PDF file, StudyHub automatically sends it to the Google Gemini AI API to generate study flashcards. The flashcards are question-and-answer pairs extracted from the document content. Students can review these flashcards using a flip-card interface in the Week Detail page.

### 6. Personal Flashcard Decks
Each student can build their own personal deck for any week by:
- Adding AI-generated flashcards to their deck
- Creating custom flashcards with their own question and answer
- Editing their copies of flashcards
- Recording their study progress (review count, difficulty rating on a 1–5 scale)
- Viewing deck statistics (total cards, reviewed cards, average difficulty, last reviewed date)

### 7. Comments
Students can post comments on the resources of any week. Comments support upvotes and downvotes. The comment feed shows who posted each comment and when.

### 8. User Profiles
Every user has a public profile showing:
- Their name and email
- All resources they have uploaded, with module context (module name, semester, year, week number)
- Action buttons to download files or open links

Users can access their own profile from the sidebar ("My Profile") and view other students' profiles by clicking on their name in any resource card.

### 9. Academic Terms
Academic Terms represent semester/year periods (e.g. "Spring 2026"). They group module runs. Admins can create new academic terms and mark one as the currently active term.

### 10. RAG Chatbot Assistant
StudyHub includes an AI-powered chatbot in the bottom-right corner of every page. The chatbot uses Retrieval-Augmented Generation (RAG) with local Ollama models to answer questions about the system. It only answers StudyHub-related questions and rejects unrelated queries.

---

## User Roles and Permissions

### Regular User (Student)
- Register and log in
- View all modules, runs, weeks
- Upload file and link resources to any week
- Delete their own resources
- Generate and view flashcards
- Build personal flashcard decks
- Post and vote on comments
- View all user profiles

### Admin User
- All regular user permissions
- Create new modules
- Create new module runs
- Create and manage academic terms
- View admin dashboard with system statistics
- Access orphan object cleanup utility

---

## Authentication

StudyHub uses JWT (JSON Web Token) authentication.

### How to log in
Send a POST request to `/api/v1/auth/login` with your email and password:

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "yourpassword"
}
```

On success you receive a JWT token:
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### How to use the token
Include the token in the `Authorization` header of every subsequent request:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Tokens are stored in the browser's `localStorage` under the key `auth_token`. The frontend automatically attaches the token to every API request via an Axios interceptor.

If a request returns HTTP 401 Unauthorized, the frontend clears the token and redirects to the login page.

### How to register
Send a POST request to `/api/v1/users` with first name, last name, email, and password:
```
POST /api/v1/users
Content-Type: application/json

{
  "first_name": "Usmon",
  "last_name": "Karimov",
  "email": "usmon@example.com",
  "password": "securepassword"
}
```

---

## API Endpoints

All endpoints are under the base path `/api/v1`. Protected endpoints require `Authorization: Bearer <token>`.

### Authentication
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/login` | No | Login and receive JWT token |
| POST | `/users` | No | Register a new user |
| GET | `/health` | No | Health check |

### Users
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/users/me` | Yes | Get current logged-in user |
| GET | `/users/{id}` | Yes | Get user by ID |
| GET | `/users` | Yes | List all users |
| DELETE | `/users/{id}` | Yes | Delete user |

### Modules
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/modules` | Yes | List all modules |
| POST | `/modules` | Yes | Create a module (admin) |
| GET | `/modules/{id}` | Yes | Get module with all runs and weeks |
| DELETE | `/modules/{id}` | Yes | Delete module |

### Module Runs
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/modules/{moduleID}/runs` | Yes | List runs for a module |
| POST | `/modules/{moduleID}/runs` | Yes | Create a module run |
| GET | `/module-runs/{id}` | Yes | Get a specific module run |
| DELETE | `/module-runs/{id}` | Yes | Delete a module run |

### Resources
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/resources/file/{week_id}` | Yes | Upload a file (multipart/form-data) |
| POST | `/resources/link/{week_id}` | Yes | Create a link resource |
| GET | `/resources/weeks/{week_id}` | Yes | List resources for a week |
| GET | `/resources/users/{user_id}` | Yes | List resources uploaded by a user |
| GET | `/resources/{id}` | Yes | Get a presigned S3 download URL |
| DELETE | `/resources/{id}` | Yes | Delete a resource |

### Comments
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/comments` | Yes | Post a comment on a week |
| GET | `/comments/weeks/{week_id}` | Yes | List comments for a week |
| POST | `/comments/{id}/upvote` | Yes | Upvote a comment |
| POST | `/comments/{id}/downvote` | Yes | Downvote a comment |

### Flashcards and Decks
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/conents/objects` | Yes | Generate flashcards from a PDF |
| POST | `/decks/weeks/{week_id}/cards` | Yes | Add an AI card to personal deck |
| POST | `/decks/weeks/{week_id}/cards/custom` | Yes | Create a custom card |
| GET | `/decks/weeks/{week_id}/cards` | Yes | Get your deck for a week |
| PATCH | `/decks/cards/{card_id}` | Yes | Edit a card |
| DELETE | `/decks/cards/{card_id}` | Yes | Remove a card from deck |
| POST | `/decks/cards/{card_id}/review` | Yes | Record a study review |
| GET | `/decks/weeks/{week_id}/stats` | Yes | Get deck statistics |

### Academic Terms
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/academic-terms/current` | Yes | Get the active academic term |
| POST | `/academic-terms/new-term` | Yes | Create a new term (admin) |

### Chatbot
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/chat` | Yes | Send a message to the RAG chatbot |

#### Chatbot request format
```
POST /api/v1/chat
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "How do I upload a file?",
  "history": [
    { "role": "user", "content": "What is StudyHub?" },
    { "role": "assistant", "content": "StudyHub is an academic study management platform..." }
  ]
}
```

#### Chatbot response format
```json
{
  "data": {
    "reply": "To upload a file, navigate to a Week Detail page and click the Upload button...",
    "sources": [
      { "source": "studyhub_docs.md", "page": null }
    ]
  }
}
```

---

## How to Use StudyHub — Step by Step

### Uploading a File
1. Log in to StudyHub
2. Go to Modules from the sidebar
3. Click on a module to open its detail page
4. Click on a module run, then select a week
5. On the Week Detail page, click **Upload File**
6. Select a file from your device and give it a name
7. Click Upload — the file is stored in S3 and flashcards are generated in the background

### Creating a Link Resource
1. Navigate to a Week Detail page
2. Click **Add Link**
3. Enter the URL and a display name
4. Click Save

### Viewing and Copying Flashcards to Your Deck
1. Navigate to a Week Detail page
2. Click the **Flashcards** tab
3. AI-generated cards appear as flippable cards
4. Click **Add to My Deck** on any card to copy it to your personal study deck
5. Click **Create Custom Card** to write your own question and answer

### Studying from Your Deck
1. Navigate to a Week Detail page
2. Click the **My Deck** tab
3. Flip through your cards
4. After each card, record a difficulty rating (1–5)
5. View your progress in the deck statistics panel

### Posting a Comment
1. Navigate to a Week Detail page
2. Scroll to the Comments section
3. Type your comment and click Post
4. Upvote or downvote other students' comments

---

## System Architecture

```
Frontend (React + TypeScript + Vite + Tailwind CSS)
    │
    │ HTTP /api/v1/*
    ▼
Backend (Go + Chi router)   ──── PostgreSQL (data storage)
    │                        ──── AWS S3 (file storage)
    │                        ──── RabbitMQ (async queue)
    │                        ──── Google Gemini API (flashcards)
    │                        ──── RAG Service (chatbot)
    ▼
RAG Service (Python + FastAPI + LangChain + ChromaDB + Ollama)
```

### Technology Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 18, TypeScript, Vite, Tailwind CSS, shadcn/ui |
| Backend | Go (Golang), Chi router, pgx (PostgreSQL driver) |
| Database | PostgreSQL 16 |
| File Storage | AWS S3 with presigned URLs |
| Message Queue | RabbitMQ |
| AI Flashcards | Google Gemini 2.5 Flash |
| Chatbot LLM | Ollama (gemma4:e4b) — runs locally |
| Chatbot Embeddings | Ollama (nomic-embed-text-v2-moe) — runs locally |
| Vector Database | ChromaDB (persisted to disk) |
| Containerisation | Docker, Docker Compose |
| CI/CD | GitHub Actions, GHCR, AWS EC2 |

### RAG Chatbot Architecture
The chatbot is a separate Python microservice that the Go backend proxies requests to.

**Build time (one-off):**  
Documents → DirectoryLoader → RecursiveCharacterTextSplitter (chunk_size=800, overlap=150) → OllamaEmbeddings (nomic-embed-text-v2-moe) → ChromaDB (persisted to disk)

**Runtime (every query):**  
User question → embed with Ollama → ChromaDB cosine similarity search (top-4 chunks) → check relevance score → build prompt with context + history → ChatOllama (gemma4:e4b, temp=0.2) → answer + source attribution

If no relevant chunks are found (score below 0.30), the chatbot responds with "I can only answer questions about StudyHub" without calling the LLM.

---

## Deployment

### Prerequisites
- Docker and Docker Compose
- AWS account with an S3 bucket
- Google Gemini API key (for flashcard generation)
- Ollama installed with models pulled:
  - `ollama pull gemma4:e4b`
  - `ollama pull nomic-embed-text-v2-moe`

### Environment Variables (.env file)
```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=studyhub
JWT_KEY=your-jwt-secret
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
AWS_S3_BUCKET=your-bucket
AWS_DEFAULT_REGION=us-east-1
AWS_S3_URL=https://your-bucket.s3.us-east-1.amazonaws.com
RBMQ_USER=guest
RBMQ_PASS=guest
RBMQ_HOST=rabbitmq
GEMINI_API_KEY=your-gemini-key
OLLAMA_BASE_URL=http://host.docker.internal:11434
RAG_SERVICE_URL=http://rag-service:8001
```

### Running with Docker Compose
```bash
# Development
docker compose -f docker-compose.dev.yml up --build

# Production
docker compose up --build
```

### Database
StudyHub uses PostgreSQL with migrations managed by golang-migrate. Migrations run automatically when using Docker Compose. Tables: users, modules, module_runs, weeks, academic_terms, storage_objects, resources, flashcards, user_deck_cards.

---

## Navigation Structure

```
/login              — Login page (public)
/register           — Registration page (public)
/modules            — List of all modules
/modules/:id        — Module detail (runs and weeks)
/modules/:moduleId/weeks/:weekId  — Week detail (resources, flashcards, comments)
/users/:userId      — User profile (their uploads)
/academic-terms     — Academic terms management
/home               — Dashboard
```

---

## Common Questions

**Q: How do I connect to the StudyHub API?**
A: Authenticate via POST /api/v1/auth/login to get a JWT token. Include it as `Authorization: Bearer <token>` in all subsequent requests.

**Q: What file types can I upload?**
A: StudyHub accepts any file type. PDFs trigger automatic flashcard generation. Files are stored in AWS S3.

**Q: Are there duplicate files?**
A: No. StudyHub uses SHA-256 hashing to detect duplicate files. If you upload the same file twice, it is stored only once in S3.

**Q: Can I see what other students uploaded?**
A: Yes. All resources in a week are visible to all authenticated users. You can also visit any student's profile page to see all their uploads.

**Q: How are flashcards generated?**
A: When you upload a PDF, the backend sends it to the Google Gemini API which extracts key concepts and generates question-answer flashcard pairs automatically.

**Q: Can I edit AI-generated flashcards?**
A: You cannot edit the original AI-generated card, but when you add it to your personal deck you get your own copy that you can edit freely.

**Q: What happens if the chatbot does not know the answer?**
A: If the question is not about StudyHub, the chatbot responds: "I can only answer questions about StudyHub." If the question is about StudyHub but not covered in the documentation, the chatbot will say so honestly.
