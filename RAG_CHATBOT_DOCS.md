# RAG Chatbot Feature — Implementation Documentation

## Overview

StudyHub now includes an AI-powered chatbot assistant that uses **Retrieval-Augmented Generation (RAG)** to answer questions grounded in uploaded study documents. Unlike a plain LLM that can only hallucinate generic answers, the RAG chatbot retrieves the most relevant passages from a vector database built from real documents before generating a response — making answers accurate, citable, and traceable to a source.

### What users can do
- Ask questions about how to use StudyHub (modules, resources, flashcards, etc.)
- Ask questions about study material that has been indexed (lecture notes, PDFs, slides)
- See which document and page number each answer came from
- The chatbot widget is available on every page — no navigation required

---

## What is RAG?

**Retrieval-Augmented Generation** is a pattern that combines a vector database (for finding relevant facts) with a large language model (for generating a natural-language response).

```
Without RAG:
  User question  ──►  LLM  ──►  Answer
                       (may hallucinate)

With RAG:
  User question  ──►  Embedding model  ──►  Vector search  ──►  Top-K chunks
                                                                      │
                       LLM  ◄── Context (chunks) + Question ──────────┘
                        │
                        ▼
                    Grounded answer + Sources
```

The key advantage is **grounding**: the LLM is explicitly given the relevant text from real documents, so it cannot invent facts.

---

## System Architecture

### Service Map

```
┌─────────────────────────────────────────────────────────┐
│                        Docker Network                    │
│                                                         │
│  ┌──────────┐   /api/v1/*   ┌──────────┐               │
│  │ Frontend │──────────────►│ Go       │               │
│  │ (React)  │               │ Backend  │               │
│  │  :80     │◄──────────────│  :8080   │               │
│  └──────────┘               └────┬─────┘               │
│                                  │                      │
│                    POST /chat     │                      │
│                                  ▼                      │
│                          ┌──────────────┐               │
│                          │  RAG Service │               │
│                          │  (Python)    │               │
│                          │   :8001      │               │
│                          └──────┬───────┘               │
│                                 │                       │
│                    ┌────────────┴────────────┐          │
│                    ▼                         ▼          │
│             ┌────────────┐          ┌──────────────┐    │
│             │  ChromaDB  │          │  Gemini API  │    │
│             │ (local FS) │          │  (external)  │    │
│             └────────────┘          └──────────────┘    │
│                                                         │
│  ┌──────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  db  │  │ rabbitmq │  │   s3     │  │gotenberg │   │
│  └──────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────┘
```

### Fallback Strategy

The Go backend always has a safety net. If the RAG service is down (e.g., during startup or a crash), it automatically falls back to a direct Gemini call with the static StudyHub system prompt.

```
Go ChatHandler
     │
     ├─► RAG service available?
     │       YES ──► call /chat ──► return answer + sources
     │       NO  ──► log warning
     │                   │
     └───────────────────► Gemini direct call ──► return answer
```

---

## Request / Response Flow

### Step-by-step: User asks a question

```
1. User types message in ChatWidget (React)
        │
        ▼
2. POST /api/v1/chat  { "message": "What is a module run?" }
        │  (JWT token in Authorization header)
        ▼
3. Go ChatHandler — validates request, checks ragServiceURL
        │
        ▼
4. POST http://rag-service:8001/chat  { "message": "..." }
        │
        ▼
5. RAG Service — embeds the question with Gemini text-embedding-004
        │
        ▼
6. ChromaDB vector search — finds top-4 most similar document chunks
        │
        ▼
7. Build prompt:
        [System prompt]
        [Retrieved chunks as context]
        [User question]
        │
        ▼
8. Gemini 2.5 Flash Lite generates answer
        │
        ▼
9. Return { "reply": "...", "sources": [{ "source": "lecture1.pdf", "page": 3 }] }
        │
        ▼
10. Go wraps in { "data": { "reply": "...", "sources": [...] } } → 200 OK
        │
        ▼
11. Frontend Axios interceptor unwraps { "data": ... }
        │
        ▼
12. ChatWidget renders answer (Markdown) + source badges
```

---

## Component Breakdown

### 1. Python RAG Service (`rag-service/`)

The self-contained microservice responsible for all RAG logic.

```
rag-service/
├── main.py          FastAPI application, lifespan startup, REST endpoints
├── rag.py           RAGChatbot class — document loading, indexing, retrieval, generation
├── requirements.txt Python dependencies
├── Dockerfile       Python 3.12-slim image
└── data/            Drop PDFs and .txt files here to index them
```

#### `rag.py` — RAGChatbot class

| Method | Responsibility |
|---|---|
| `__init__` | Creates Gemini embeddings client and Gemini LLM client |
| `_init_db` | Loads existing ChromaDB if present, otherwise calls `rebuild()` |
| `_load_documents` | Uses LangChain `DirectoryLoader` to load all PDFs and .txt files from `data/` |
| `rebuild` | Splits documents into 800-char chunks with 150-char overlap, embeds them, saves to ChromaDB |
| `chat` | Embeds question, retrieves top-4 chunks, builds prompt, calls Gemini, returns answer + sources |

#### `main.py` — FastAPI endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Liveness check — returns `{"status": "ok"}` |
| `POST` | `/chat` | Body: `{"message": "..."}` — returns `{"reply": "...", "sources": [...]}` |
| `POST` | `/rebuild` | Re-scans `data/` and rebuilds the vector index — use after adding new documents |

#### Chunking strategy

Documents are split using `RecursiveCharacterTextSplitter` with:
- **chunk_size = 800** characters — small enough to stay focused, large enough for context
- **chunk_overlap = 150** characters — overlap prevents answers from being cut at boundaries

#### Embeddings and model

| Component | Model | Why |
|---|---|---|
| Embeddings | `models/text-embedding-004` (Google) | Same API key as existing Gemini usage, high quality |
| LLM | `gemini-2.5-flash-lite` (Google) | Fast, cheap, already used in the project |
| Vector DB | ChromaDB (local filesystem) | Zero infrastructure cost, persists between restarts |

---

### 2. Go Backend Changes (`backend/`)

Three files were modified to wire the RAG service into the existing request pipeline.

#### `internal/config/config.go`

Added one field:
```go
RAGServiceURL string `env:"RAG_SERVICE_URL" envDefault:"http://rag-service:8001"`
```

The default points at the Docker service name. For local development without Docker, set `RAG_SERVICE_URL=http://localhost:8001`.

#### `internal/http/http.go`

Added `ragServiceURL string` to `HTTPServer` struct and to `NewHTTPServer(...)` constructor. No routing changes — the `/chat` route was already registered.

#### `internal/http/chat_handler.go`

`ChatHandler` now:
1. Calls `callRAGService(srv.ragServiceURL, req.Message)` — a private function that POSTs to the Python service
2. On success: returns `chatResponse{Reply, Sources}` to the frontend
3. On any error (service down, timeout, bad status): logs a warning and falls back to `srv.geminiClient.Chat()`

The HTTP client used by `callRAGService` has a 30-second timeout, matching a reasonable LLM response time.

**Updated response struct:**
```go
type chatResponse struct {
    Reply   string       `json:"reply"`
    Sources []chatSource `json:"sources,omitempty"`  // new field
}

type chatSource struct {
    Source string `json:"source"`   // filename, e.g. "lecture3.pdf"
    Page   *int   `json:"page,omitempty"` // 0-indexed page number
}
```

---

### 3. Frontend Changes (`frontend/`)

#### `src/api/chat.ts`

The API function now returns a structured object instead of a plain string:
```typescript
export interface ChatSource {
  source: string
  page?: number
}

export interface ChatReply {
  reply: string
  sources: ChatSource[]
}
```

#### `src/components/layout/ChatWidget.tsx`

Each assistant message now stores its sources alongside the content:
```typescript
interface Message {
  role: 'user' | 'assistant'
  content: string
  sources?: ChatSource[]   // new
}
```

Source badges are rendered below each RAG-grounded reply:

```
┌──────────────────────────────────────────────┐
│  A module run is a specific instance of a    │
│  module delivered in a given semester...      │
│                                              │
│  [📄 week1-notes.pdf p.4] [📄 syllabus.pdf] │
└──────────────────────────────────────────────┘
```

---

## Data Flow Diagrams

### Document Indexing (one-time or on /rebuild)

```
data/ folder
    │
    ├── lecture1.pdf  ──┐
    ├── lecture2.pdf  ──┤  DirectoryLoader (LangChain)
    ├── notes.txt     ──┘
              │
              ▼
    RecursiveCharacterTextSplitter
    chunk_size=800, overlap=150
              │
              ▼
    [ chunk_1 ][ chunk_2 ][ chunk_3 ] ... [ chunk_N ]
              │
              ▼
    GoogleGenerativeAIEmbeddings
    (text-embedding-004)
              │
              ▼
    768-dimensional float vectors
              │
              ▼
    ChromaDB  (persisted to ./chroma_db/)
```

### Query / Answer Flow

```
User question: "What is a module run?"
              │
              ▼
    GoogleGenerativeAIEmbeddings (same model)
              │
              ▼
    768-dim query vector
              │
              ▼
    ChromaDB cosine similarity search  →  top-4 chunks
              │
              ▼
    Prompt assembly:
    ┌─────────────────────────────────────────┐
    │ [System: StudyHub assistant context]    │
    │ [Context chunk 1]                       │
    │ [Context chunk 2]                       │
    │ [Context chunk 3]                       │
    │ [Context chunk 4]                       │
    │ User: What is a module run?             │
    │ Assistant:                              │
    └─────────────────────────────────────────┘
              │
              ▼
    Gemini 2.5 Flash Lite
              │
              ▼
    "A module run is a specific instance of a
     module delivered in a given semester..."
              │
              ▼
    Sources: [{ source: "syllabus.pdf", page: 2 }]
```

---

## API Reference

### RAG Service (internal, port 8001)

#### `GET /health`
```
Response 200:
{ "status": "ok" }
```

#### `POST /chat`
```
Request:
{ "message": "string" }

Response 200:
{
  "reply": "string",
  "sources": [
    { "source": "filename.pdf", "page": 3 }
  ]
}

Response 400:
{ "detail": "message is required" }
```

#### `POST /rebuild`
```
Response 200 (documents found):
{
  "status": "ok",
  "documents": 5,
  "chunks": 142,
  "message": "Indexed 5 document(s) into 142 chunks."
}

Response 200 (no documents):
{
  "status": "empty",
  "chunks": 0,
  "message": "No documents found. Drop PDFs or .txt files into data/ and call /rebuild."
}
```

### Go Backend (public, port 8080)

#### `POST /api/v1/chat`

Requires `Authorization: Bearer <jwt_token>`.

```
Request:
{ "message": "string" }

Response 200:
{
  "data": {
    "reply": "string",
    "sources": [
      { "source": "lecture1.pdf", "page": 2 }
    ]
  }
}

Response 400:
{ "error": { "code": 400, "message": "message is required" } }

Response 500:
{ "error": { "code": 500, "message": "failed to get response from AI" } }
```

---

## Setup and Configuration

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `GEMINI_API_KEY` | *(required)* | Google Gemini API key — used by both the RAG service and the Go backend |
| `RAG_SERVICE_URL` | `http://rag-service:8001` | URL the Go backend uses to reach the RAG service |
| `DATA_DIR` | `data` | Directory the RAG service scans for documents |
| `DB_DIR` | `chroma_db` | Directory where ChromaDB persists its vector index |

### Adding Documents to the Knowledge Base

1. Place `.pdf` or `.txt` files into `rag-service/data/`
2. Trigger a re-index:

```bash
# In development (port is exposed):
curl -X POST http://localhost:8001/rebuild

# Via docker exec:
docker compose exec rag-service curl -X POST http://localhost:8001/rebuild
```

3. The vector database is rebuilt in-place — no restart needed.

### Running with Docker Compose

**Development** (source code mounted, RAG service port exposed):
```bash
docker compose -f docker-compose.dev.yml up --build
```
RAG service admin endpoints are accessible at `http://localhost:8001`.

**Production:**
```bash
docker compose up --build
```
RAG service is internal-only (no published port).

### Running the RAG Service Locally (without Docker)

```bash
cd rag-service
pip install -r requirements.txt
export GEMINI_API_KEY=your-key-here
uvicorn main:app --host 0.0.0.0 --port 8001 --reload
```

Then set `RAG_SERVICE_URL=http://localhost:8001` for the Go backend.

---

## Technology Choices

| Choice | Alternative considered | Reason |
|---|---|---|
| **Gemini embeddings** | Ollama (local) | No extra installation needed; same API key already in the project |
| **ChromaDB** | Pinecone, Weaviate | Zero cost, runs locally, persists to disk, no extra service |
| **FastAPI** | Flask, Django | Async-ready, automatic OpenAPI docs, Pydantic validation |
| **Python microservice** | Embedding RAG into Go | LangChain ecosystem is Python-native; avoids porting complex logic |
| **Gemini 2.5 Flash Lite** | GPT-4, Claude | Already in use, lowest cost, fast response time |
| **Fallback to direct Gemini** | Return error if RAG is down | Guarantees chatbot always works even if RAG service crashes |

---

## Limitations and Future Work

### Current Limitations

- **Manual re-indexing** — documents must be added to `data/` manually and `/rebuild` called; there is no automatic sync with S3 resources
- **Single language** — ChromaDB and the embedding model work best with English documents
- **No authentication on RAG endpoints** — the `/rebuild` endpoint is unauthenticated; in production it is internal-only (no exposed port)
- **Cold start** — if ChromaDB is empty, the first startup builds the index, which takes time proportional to document count

### Potential Enhancements

1. **Auto-index S3 resources** — after a file upload, push the object key to a queue; a worker downloads and adds it to ChromaDB automatically
2. **Per-module knowledge bases** — separate vector collections per module so retrieval is scoped to the relevant course
3. **Admin rebuild UI** — add a button in the admin dashboard to trigger `/rebuild` without needing `curl`
4. **Streaming responses** — use FastAPI `StreamingResponse` and server-sent events so the answer appears word-by-word
5. **Conversation history** — pass the last N message pairs to Gemini for multi-turn context

---

## Files Changed / Created

### New Files

| Path | Description |
|---|---|
| `rag-service/main.py` | FastAPI application with `/health`, `/chat`, `/rebuild` endpoints |
| `rag-service/rag.py` | `RAGChatbot` class — document loading, ChromaDB indexing, retrieval, generation |
| `rag-service/requirements.txt` | Python dependencies |
| `rag-service/Dockerfile` | Python 3.12-slim Docker image |
| `rag-service/data/.gitkeep` | Placeholder for the document input directory |
| `rag-service/.gitignore` | Ignores `chroma_db/` and `__pycache__/` |

### Modified Files

| Path | Change |
|---|---|
| `backend/internal/config/config.go` | Added `RAGServiceURL` field |
| `backend/internal/http/http.go` | Added `ragServiceURL` to `HTTPServer` struct and constructor |
| `backend/internal/http/chat_handler.go` | Proxies to RAG service; falls back to Gemini |
| `backend/cmd/main.go` | Passes `cfg.RAGServiceURL` to `NewHTTPServer` |
| `docker-compose.yml` | Added `rag-service` with persistent volumes |
| `docker-compose.dev.yml` | Added `rag-service` with local volume mount and exposed port |
| `frontend/src/api/chat.ts` | Exports `ChatSource`, `ChatReply` types; returns structured object |
| `frontend/src/components/layout/ChatWidget.tsx` | Displays source badges below assistant messages |
