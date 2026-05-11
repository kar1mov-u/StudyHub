# RAG-Based Chatbot — Technical Documentation

**Module:** 4-SEI-10-899 | Software Engineering and Internship  
**Team Project:** StudyHub  
**Spring 2025-2026 | Central Asian University**

---

## Overview

For this assessment we integrated a RAG-based chatbot into our existing StudyHub application. The chatbot acts as a help-desk assistant that can only answer questions about StudyHub — its features, APIs, authentication, and architecture. It is built on top of fully local models running through Ollama, so no data is sent to any external API at chat time.

The main things the chatbot can do:
- Answer questions about StudyHub using retrieved context from our documentation
- Reject any question that is not about StudyHub with a consistent message
- Remember the last few turns of a conversation so follow-up questions work
- Show the user which document the answer came from (source attribution)

---

## What is RAG and Why We Used It

RAG (Retrieval-Augmented Generation) solves a core problem with LLMs: they hallucinate. If you just ask a plain LLM about StudyHub it will make things up, because it has never seen our documentation.

With RAG, instead of relying on the model's internal memory, we first search our own documentation for relevant chunks, then pass those chunks to the model as context. The model can only answer from what we give it.

```
Without RAG:
  User question ──► LLM ──► Made-up answer

With RAG:
  User question ──► embed ──► search ChromaDB ──► top-4 relevant chunks
                                                          │
                          LLM ◄── system prompt + chunks ─┘
                           │
                           ▼
                    Grounded answer + source attribution
```

If no relevant chunk is found above the similarity threshold, we reject the question before the LLM is even called. This makes rejection fast and consistent.

---

## Architecture

### Service Map

```
┌──────────────────────────────────────────────────────────────┐
│                        Docker Network                         │
│                                                              │
│  ┌──────────┐   /api/v1/*    ┌────────────┐                  │
│  │ Frontend │───────────────►│ Go Backend │                  │
│  │  React   │◄───────────────│   :8080    │                  │
│  │  :80     │                └─────┬──────┘                  │
│  └──────────┘                      │ POST /chat              │
│                                    ▼                         │
│                          ┌──────────────────┐                │
│                          │   RAG Service    │                │
│                          │  Python/FastAPI  │                │
│                          │     :8001        │                │
│                          └────────┬─────────┘               │
│                                   │                          │
│                    ┌──────────────┴──────────────┐           │
│                    ▼                             ▼           │
│             ┌────────────┐             ┌──────────────┐      │
│             │  ChromaDB  │             │    Ollama    │      │
│             │  (disk)    │             │  (host)      │      │
│             └────────────┘             └──────────────┘      │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

The RAG service is a separate Python microservice. We kept it separate from the Go backend because the LangChain ecosystem is Python-native, and trying to do RAG in Go would have been significantly more complex. The Go backend proxies chat requests to it and falls back to a direct Gemini call if the RAG service is unavailable.

### Models

We used two local Ollama models as recommended in the assessment spec:

| Role | Model | Size | Purpose |
|------|-------|------|---------|
| Embeddings | `nomic-embed-text-v2-moe` | 957 MB | Converts text to vectors for similarity search |
| LLM | `gemma4:e4b` | ~3 GB | Reads retrieved chunks and generates the answer |

Both models run on the host machine through Ollama. The Docker containers connect to them via `host.docker.internal:11434`.

### Fallback

If the RAG service or Ollama is down, the Go backend falls back to a direct Gemini API call so the chat widget does not break completely. The fallback does not have source attribution.

```
ChatHandler receives request
        │
        ├─► RAG service available?
        │       YES ──► returns grounded answer + sources
        │       NO  ──► logs warning
        │                    │
        └────────────────────► Gemini direct call ──► answer (no sources)
```

---

## Build Time vs. Runtime Separation

This was one of the most important design decisions we made. Document ingestion (loading files, chunking, embedding) is expensive — it calls the embedding model for every chunk. Doing this on every startup or on every user question would be completely impractical.

Our solution: run ingestion exactly once at build time, persist the results to disk, and just load from disk on every subsequent start.

```
BUILD TIME (runs once)                  RUNTIME (runs on every question)
──────────────────────                  ────────────────────────────────
data/studyhub_docs.md                   User question
         │                                      │
         ▼                                      ▼
  DirectoryLoader                       OllamaEmbeddings
  (PDF, TXT, MD)                        (nomic-embed-text-v2-moe)
         │                                      │
         ▼                                      ▼
  RecursiveCharacterTextSplitter        ChromaDB similarity search
  chunk_size=800, overlap=150           top-4 chunks + cosine scores
         │                                      │
         ▼                               relevance check (≥ 0.30)
  OllamaEmbeddings                             │
  (nomic-embed-text-v2-moe)             ┌──────┴──────┐
         │                            PASS           FAIL
         ▼                              │               │
  ChromaDB                        build prompt    "I can only answer
  (persisted to chroma_db/)             │         questions about StudyHub"
                                        ▼
                                 ChatOllama (gemma4:e4b, temp=0.2)
                                        │
                                        ▼
                                 answer + sources
```

The build step is handled by `build_index.py`. When the Docker container starts, it checks if `chroma_db/` already has data. If yes, it exits immediately and the server starts in under a second. If not (first start), it builds the index. This means `docker compose up` is always fast after the first run.

```python
# build_index.py — skip if index already exists
if not force and os.path.exists(DB_DIR) and any(os.scandir(DB_DIR)):
    print("Vector index already exists — skipping build.")
    sys.exit(0)

# First run: build the full index
bot = RAGChatbot()
result = bot.rebuild()
# Output: "Indexed 1 document(s) into 25 chunks."
```

We also expose `POST /rebuild` as an endpoint so we can re-index without restarting the container after adding new documents.

---

## Implementation Details

### File Structure

```
rag-service/
├── main.py           FastAPI app — endpoints, request/response models
├── rag.py            RAGChatbot class — all RAG logic
├── build_index.py    Build-time script — indexes data/ into ChromaDB
├── requirements.txt  Python dependencies
├── Dockerfile        Python 3.12-slim, runs build_index.py then uvicorn
└── data/
    └── studyhub_docs.md   Knowledge base (~600 lines of StudyHub documentation)
```

### Initialisation (`rag.py`)

```python
from langchain_ollama import ChatOllama, OllamaEmbeddings

class RAGChatbot:
    def __init__(self):
        self.embeddings = OllamaEmbeddings(
            model="nomic-embed-text-v2-moe",
            base_url="http://localhost:11434",
        )
        self.llm = ChatOllama(
            model="gemma4:e4b",
            base_url="http://localhost:11434",
            temperature=0.2,  # low temperature = factual, not creative
        )
        self.vector_db = None
        self._init_db()  # load existing index, or build fresh if it doesn't exist
```

### Chunking and Indexing

```python
def rebuild(self):
    docs = self._load_documents()  # reads all .pdf, .txt, .md files from data/

    splitter = RecursiveCharacterTextSplitter(
        chunk_size=800,    # roughly one paragraph per chunk
        chunk_overlap=150  # 150-char overlap so sentences on chunk boundaries aren't lost
    )
    chunks = splitter.split_documents(docs)

    self.vector_db = Chroma.from_documents(
        documents=chunks,
        embedding=self.embeddings,
        persist_directory="chroma_db",  # saved to disk — only runs once
    )
```

We chose `chunk_size=800` with `chunk_overlap=150` because it keeps each chunk focused on a single topic while ensuring that content spanning a boundary is not lost. The 800-character limit is roughly one paragraph, which is a natural unit of meaning for a help document.

### Retrieval and Rejection (FR-01, FR-03)

```python
# Retrieve top-4 chunks with cosine similarity scores
results = self.vector_db.similarity_search_with_relevance_scores(
    retrieval_query, k=4
)

# Filter by relevance threshold
relevant = [(doc, score) for doc, score in results if score >= 0.30]

if not relevant and results:
    # All chunks scored below threshold — question is off-topic
    return "I can only answer questions about StudyHub.", []
```

The threshold of **0.30** acts as a hard gate. If no chunk from our documentation scores at least 0.30 cosine similarity with the question, the system rejects it immediately without calling the LLM at all. This is efficient (no token cost) and consistent (the rejection message is always the same).

### Multi-Turn Context (FR-06)

One issue we ran into: for short follow-up questions like "where is it stored?", the question alone is too vague for vector search to find anything relevant. The embedding of that phrase doesn't match anything specific in the documentation.

Our solution was to enrich the retrieval query by prepending the previous user message. The LLM still receives only the current question, but the retrieval benefits from the added context.

```python
def chat(self, question: str, history: list = None):
    # Enrich retrieval query with previous user message for follow-up questions
    retrieval_query = question
    if history:
        last_user = next(
            (m["content"] for m in reversed(history) if m.get("role") == "user"),
            None
        )
        if last_user:
            retrieval_query = f"{last_user} {question}"

    # ... retrieval ...

    # Build message list with conversation history for multi-turn context
    messages = [SystemMessage(content=system_with_context)]
    for msg in (history or [])[-6:]:  # last 6 messages = 3 full exchanges
        if msg["role"] == "user":
            messages.append(HumanMessage(content=msg["content"]))
        elif msg["role"] == "assistant":
            messages.append(AIMessage(content=msg["content"]))

    messages.append(HumanMessage(content=question))
    response = self.llm.invoke(messages)
```

---

## Prompt Engineering Strategy

The system prompt is injected at the top of every LLM request. We designed it with two goals: tell the model its role, and give it explicit rules for rejection.

```python
_SYSTEM = """You are StudyHub Assistant, the official help-desk chatbot for the StudyHub application.

StudyHub is an academic study management platform. Students use it to organise modules, share study
resources (files and links), generate AI-powered flashcards, and collaborate via comments.

YOUR STRICT RULES:
1. Only answer questions about StudyHub — its features, APIs, authentication, architecture, and usage.
2. Base every answer on the CONTEXT PROVIDED below. Do not invent facts.
3. If the provided context does not contain enough information to answer confidently, respond with exactly:
   "I can only answer questions about StudyHub."
4. If the question is clearly unrelated to StudyHub (e.g. history, maths, other software, personal advice),
   respond with exactly: "I can only answer questions about StudyHub."
5. Be concise. Use bullet points or numbered lists where helpful.
"""
```

A few decisions we made here:

**Why give the exact rejection phrase in the prompt?** We want the rejection to be consistent and predictable. If we just say "reject unrelated questions politely", different phrasings would appear. By specifying the exact string, the UI can also detect it reliably if needed.

**Why temperature=0.2?** A higher temperature makes the model more creative and varied — good for creative writing, bad for a help desk. We want factual, reproducible answers. Low temperature keeps the model close to what the retrieved context says.

**Two-layer rejection:** We deliberately have rejection at two levels:
1. **Score gate** (before LLM): if no chunk scores ≥ 0.30, return rejection immediately
2. **Prompt instruction** (inside LLM): even if chunks are retrieved, tell the model to reject if context is insufficient

This means even a borderline retrieval that sneaks past the score gate still gets rejected if the LLM judges the context insufficient.

---

## Knowledge Base

The chatbot's knowledge comes entirely from `rag-service/data/studyhub_docs.md` — a document we wrote covering everything about StudyHub. We chose Markdown over PDF because it is easier to maintain and update.

| Section | What it covers |
|---------|----------------|
| System overview | What StudyHub is, who it's for, the core value proposition |
| Core features | Modules, weeks, resources (files/links), flashcards, decks, comments, academic terms |
| Authentication | JWT login flow, where the token is stored, how to send it in requests |
| API endpoints | All 29 endpoints — method, path, auth requirement, request/response examples |
| How-to guides | Step-by-step instructions for every main workflow |
| Architecture | Tech stack, deployment diagram, how the services connect |
| User roles | Admin vs. regular user permissions |
| Common Q&A | Pre-written answers to the most likely questions |

This document is baked into the Docker image at build time. On first container start it gets chunked, embedded, and stored in ChromaDB. The index then persists across restarts via a named Docker volume.

---

## Integration into the Application

### Backend (Go)

The existing Go backend has a `ChatHandler` that forwards requests to the RAG service. We added conversation history support — the handler accepts a `history` array and passes it through:

```go
type chatRequest struct {
    Message string               `json:"message"`
    History []chatHistoryMessage `json:"history"`
}

func (srv *HTTPServer) ChatHandler(w http.ResponseWriter, r *http.Request) {
    var req chatRequest
    json.NewDecoder(r.Body).Decode(&req)

    if srv.ragServiceURL != "" {
        if reply, sources, err := callRAGService(srv.ragServiceURL, req.Message, req.History); err == nil {
            ResponseWithJSON(w, http.StatusOK, chatResponse{Reply: reply, Sources: sources})
            return
        }
        slog.Warn("rag service unavailable, falling back to gemini")
    }

    reply, _ := srv.geminiClient.Chat(r.Context(), req.Message)
    ResponseWithJSON(w, http.StatusOK, chatResponse{Reply: reply})
}
```

The HTTP client timeout is set to 60 seconds because local Ollama inference on first call can be slow while the model loads into memory.

### Frontend (React)

The chat widget sends the full conversation history with every message. Source attribution is displayed as small file badges below each bot reply.

```typescript
const send = async () => {
    const text = input.trim()
    setMessages(prev => [...prev, { role: 'user', content: text }])

    const history = messages.map(m => ({ role: m.role, content: m.content }))
    const { reply, sources } = await chatApi.sendMessage(text, history)

    setMessages(prev => [...prev, { role: 'assistant', content: reply, sources }])
}
```

```tsx
{msg.sources && msg.sources.length > 0 && (
    <div className="flex flex-wrap gap-1 px-1">
        {msg.sources.map((s, j) => (
            <span key={j} className="inline-flex items-center gap-1 text-[10px]
                                     text-gray-500 bg-gray-100 rounded px-1.5 py-0.5">
                <FileText className="h-2.5 w-2.5" />
                {s.source}{s.page != null ? ` p.${s.page + 1}` : ''}
            </span>
        ))}
    </div>
)}
```

---

## Deployment

### Prerequisites

```bash
# Install Ollama then pull the two models (one-time setup)
ollama pull nomic-embed-text-v2-moe   # 957 MB
ollama pull gemma4:e4b                 # ~3 GB
```

### Running locally (no Docker)

```bash
cd rag-service
pip install -r requirements.txt

# Build the vector index (runs once, skips on subsequent runs)
python3 build_index.py

# Start the service
OLLAMA_BASE_URL=http://localhost:11434 uvicorn main:app --port 8001 --reload
```

Set `RAG_SERVICE_URL=http://localhost:8001` in the Go backend env.

### Running with Docker Compose

```bash
# Development (RAG port exposed)
docker compose -f docker-compose.dev.yml up --build

# Production
docker compose up --build
```

The `OLLAMA_BASE_URL` in both compose files is set to `http://host.docker.internal:11434` so the container can reach the Ollama instance running on the host machine.

### Dockerfile

```dockerfile
FROM python:3.12-slim
WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY main.py rag.py build_index.py ./
COPY data/ ./data/

RUN mkdir -p chroma_db
EXPOSE 8001

# On first start: build index then serve.
# On subsequent starts: index already exists, skip straight to uvicorn.
CMD ["sh", "-c", "python3 build_index.py && uvicorn main:app --host 0.0.0.0 --port 8001"]
```

### Rebuilding the index after adding documents

```bash
# Drop new files into rag-service/data/ then call the rebuild endpoint
curl -X POST http://localhost:8001/rebuild

# Or force rebuild from inside the container
docker compose exec rag-service python3 build_index.py --force
```

---

## Test Evidence

### System Q&A (10 cases)

| # | Question | Bot Reply (summarised) | Source |
|---|----------|------------------------|--------|
| 1 | What is StudyHub? | StudyHub is an academic study management platform where students organise modules, share resources, and generate flashcards. | studyhub_docs.md |
| 2 | How do I log in? | Send a POST request to `/api/v1/auth/login` with your email and password. The response includes a JWT token to use in the Authorization header. | studyhub_docs.md |
| 3 | Where is the auth token stored? | The JWT token is stored in `localStorage` under the key `auth_token`. It is automatically attached to every API request by the Axios interceptor. | studyhub_docs.md |
| 4 | How do I upload a file? | Navigate to a Week Detail page, click Upload File, select your file, and confirm. The file is sent as multipart/form-data to `POST /api/v1/resources/file/{week_id}`. | studyhub_docs.md |
| 5 | What is a module run? | A module run is a specific instance of a module in a given semester and year. One module can have multiple runs — for example, the same course in Spring 2025 and Fall 2025. | studyhub_docs.md |
| 6 | Can I add a link as a resource? | Yes. On a Week Detail page click Add Link, provide a name and URL, and submit. This calls `POST /api/v1/resources/link/{week_id}`. | studyhub_docs.md |
| 7 | What API endpoint lists all modules? | `GET /api/v1/modules` returns a list of all available modules. No auth is required. | studyhub_docs.md |
| 8 | How does the flashcard feature work? | You can generate AI-powered flashcards from your study materials. Flashcards are grouped into decks and can be reviewed in the application. | studyhub_docs.md |
| 9 | What is the tech stack? | The backend is Go with the Chi router and PostgreSQL. The frontend is React + TypeScript with Tailwind CSS. Files are stored on AWS S3. The chatbot runs on Python/FastAPI with Ollama. | studyhub_docs.md |
| 10 | What is the difference between admin and regular users? | Admin users can manage modules, academic terms, and all system content. Regular users can upload resources to weeks they have access to and manage their own content. | studyhub_docs.md |

### Rejection Cases (5 cases)

| # | Question | Bot Reply |
|---|----------|-----------|
| 1 | What is the capital of France? | I can only answer questions about StudyHub. |
| 2 | Can you write me a Python script to sort a list? | I can only answer questions about StudyHub. |
| 3 | Who won the 2022 World Cup? | I can only answer questions about StudyHub. |
| 4 | What do you think about AI ethics? | I can only answer questions about StudyHub. |
| 5 | Can you explain how blockchain works? | I can only answer questions about StudyHub. |

### Multi-Turn Test

| Turn | Speaker | Message |
|------|---------|---------|
| 1 | User | How do I authenticate with the API? |
| 1 | Bot | Send a POST request to `/api/v1/auth/login` with your email and password. You will receive a JWT token in the response. Include it in all subsequent requests as `Authorization: Bearer <token>`. |
| 2 | User | Where is the token stored in the frontend? |
| 2 | Bot | The token is stored in `localStorage` under the key `auth_token`. The Axios client in `frontend/src/api/client.ts` reads it automatically and adds the Authorization header to every request. |
| 3 | User | What happens if it expires? |
| 3 | Bot | The Axios response interceptor detects a 401 response and redirects the user to the login page. They must log in again to get a new token. |

---

## Functional Requirements Coverage

| Requirement | How it is satisfied |
|-------------|---------------------|
| **FR-01** Context-aware responses | ChromaDB similarity search retrieves top-4 relevant chunks before every LLM call |
| **FR-02** System knowledge coverage | `studyhub_docs.md` covers features, all 29 API endpoints, authentication, deployment, architecture, and user roles |
| **FR-03** Rejection of unrelated queries | Hard rejection via cosine score gate (< 0.30) + LLM system prompt instruction |
| **FR-04** Conversational interface | ChatWidget in React shows message thread, typing indicator, and distinct user/bot styling |
| **FR-05** Source attribution | Every answer returns a `sources` array rendered as file badge chips below the reply |
| **FR-06** Multi-turn context | Last 6 messages sent as history; retrieval query enriched with previous user turn for follow-ups |
| **Build/Runtime separation** | `build_index.py` runs once at container start; ChromaDB index persists to disk; `_init_db()` just opens the existing index on subsequent starts |
