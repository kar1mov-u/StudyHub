import os
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from rag import RAGChatbot

bot: RAGChatbot | None = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global bot
    api_key = os.environ.get("GEMINI_API_KEY", "")
    if not api_key:
        raise RuntimeError("GEMINI_API_KEY environment variable is required")
    bot = RAGChatbot(api_key)
    yield


app = FastAPI(title="StudyHub RAG Service", lifespan=lifespan)


class ChatRequest(BaseModel):
    message: str


class Source(BaseModel):
    source: str
    page: int | None = None


class ChatResponse(BaseModel):
    reply: str
    sources: list[Source] = []


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/chat", response_model=ChatResponse)
def chat(req: ChatRequest):
    if not req.message.strip():
        raise HTTPException(status_code=400, detail="message is required")
    reply, sources = bot.chat(req.message)
    return ChatResponse(reply=reply, sources=[Source(**s) for s in sources])


@app.post("/rebuild")
def rebuild():
    return bot.rebuild()
