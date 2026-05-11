from contextlib import asynccontextmanager
from typing import List, Optional

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from rag import RAGChatbot

bot = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Start the RAG chatbot on application startup (loads ChromaDB from disk)."""
    global bot
    bot = RAGChatbot()
    yield


app = FastAPI(title="StudyHub RAG Service", lifespan=lifespan)


# ── Request / Response models ─────────────────────────────────────────────────

class HistoryMessage(BaseModel):
    role: str      # "user" or "assistant"
    content: str


class ChatRequest(BaseModel):
    message: str
    history: List[HistoryMessage] = []   # previous turns for multi-turn context (FR-06)


class Source(BaseModel):
    source: str
    page: Optional[int] = None


class ChatResponse(BaseModel):
    reply: str
    sources: List[Source] = []


# ── Endpoints ─────────────────────────────────────────────────────────────────

@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/chat", response_model=ChatResponse)
def chat(req: ChatRequest):
    """Answer a question using RAG. Rejects off-topic questions (FR-03)."""
    if not req.message.strip():
        raise HTTPException(status_code=400, detail="message is required")

    history = [{"role": m.role, "content": m.content} for m in req.history]
    reply, sources = bot.chat(req.message, history)

    return ChatResponse(reply=reply, sources=[Source(**s) for s in sources])


@app.post("/rebuild")
def rebuild():
    """
    Re-scan the data/ folder and rebuild the vector index.
    Call this after adding new PDF/TXT/MD files to data/.
    This is the BUILD-TIME step — expensive but only needed once per document set.
    """
    return bot.rebuild()
