import os
from typing import Optional

from langchain_ollama import ChatOllama, OllamaEmbeddings
from langchain_chroma import Chroma
from langchain_community.document_loaders import PyPDFLoader, TextLoader, DirectoryLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.messages import HumanMessage, AIMessage, SystemMessage

# ── Configuration ────────────────────────────────────────────────────────────
# All values can be overridden via environment variables.

# Where Ollama is listening. On Mac/Windows with Docker use host.docker.internal;
# locally (no Docker) use localhost.
OLLAMA_BASE_URL = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")

# Local LLM for generating answers — lightweight 4-bit Gemma 4 (professor's spec).
LLM_MODEL = os.getenv("LLM_MODEL", "gemma4:e4b")

# Local embedding model — converts text to vectors for similarity search.
EMBED_MODEL = os.getenv("EMBED_MODEL", "nomic-embed-text-v2-moe")

DATA_DIR = os.getenv("DATA_DIR", "data")      # folder with PDF / TXT knowledge-base files
DB_DIR   = os.getenv("DB_DIR",   "chroma_db") # where ChromaDB persists its vector index

# Cosine-similarity threshold: a retrieved chunk scoring below this is considered
# off-topic. Lower = more permissive; higher = stricter rejection.
RELEVANCE_THRESHOLD = float(os.getenv("RELEVANCE_THRESHOLD", "0.30"))

# ── System prompt ─────────────────────────────────────────────────────────────
# This is injected at the top of every LLM request. It tells the model its role
# and — critically — instructs it to reject anything not about StudyHub.
_SYSTEM = """You are StudyHub Assistant, the help-desk chatbot for the StudyHub application.

StudyHub is an academic study management platform. Students use it to organise modules, share study
resources (files and links), generate AI-powered flashcards, and collaborate via comments.

Here is what users can do in StudyHub:
- Browse and join Modules (courses) organised by department
- Each module has one or more Runs (a specific semester and year instance)
- Each Run has weekly sessions — Weeks — that hold uploaded files, links, and comments
- Upload files (PDFs, Word docs, PowerPoints) or add external links to any week
- AI flashcard decks are automatically generated from uploaded PDFs in the background
- Comment on weekly resources and upvote or downvote others' comments
- View user profiles to see all resources a student has uploaded with full module context
- Navigate: Modules → Module Detail → Week Detail → Resources / Comments / Flashcards

YOUR RULES:
1. Answer questions about StudyHub — its features, APIs, authentication, architecture, and usage.
2. Respond to greetings and meta-questions (who you are, what you can help with, what the app does).
3. Use the CONTEXT PROVIDED below when it is relevant. For general platform questions, answer from
   the feature knowledge above even if no context was retrieved.
4. Only refuse if the question is clearly unrelated to StudyHub or academic study tools
   (e.g. sports results, geography, politics, recipes). Say exactly: "I can only answer questions about StudyHub."
5. Be friendly and concise. Use bullet points where helpful.

"""

# Questions about the platform itself bypass the score gate — they go straight to the LLM
# which can answer them from the system prompt above.
_PLATFORM_SIGNALS = {
    "app", "platform", "studyhub", "what can", "can you", "who are you",
    "what do you", "what does", "how do i", "how does", "feature", "help me",
    "introduce", "tell me about", "explain",
}


class RAGChatbot:
    def __init__(self):
        # Embedding model — turns text into a list of numbers so we can do
        # semantic (meaning-based) similarity search instead of keyword matching.
        self.embeddings = OllamaEmbeddings(
            model=EMBED_MODEL,
            base_url=OLLAMA_BASE_URL,
        )

        # LLM — reads the retrieved context + the user's question and writes an answer.
        # temperature=0.2 keeps answers factual and deterministic (not creative/random).
        self.llm = ChatOllama(
            model=LLM_MODEL,
            base_url=OLLAMA_BASE_URL,
            temperature=0.2,
        )

        self.vector_db = None
        self._init_db()

    # ── Build-time helpers ────────────────────────────────────────────────────

    def _init_db(self):
        """Load the existing ChromaDB index from disk, or build it fresh."""
        # ChromaDB persists its vectors to chroma_db/ on disk.
        # On startup we just open that folder — no re-embedding needed.
        # Only on the very first run (or after /rebuild) is the index created.
        if os.path.exists(DB_DIR) and any(os.scandir(DB_DIR)):
            self.vector_db = Chroma(
                persist_directory=DB_DIR,
                embedding_function=self.embeddings,
            )
        else:
            self.rebuild()

    def _load_documents(self):
        """Read every PDF and .txt file from the data/ knowledge-base folder."""
        os.makedirs(DATA_DIR, exist_ok=True)
        docs = []
        for loader_cls, glob in [(PyPDFLoader, "**/*.pdf"), (TextLoader, "**/*.txt"),
                                  (TextLoader, "**/*.md")]:
            try:
                loader = DirectoryLoader(
                    DATA_DIR,
                    glob=glob,
                    loader_cls=loader_cls,
                    silent_errors=True,
                )
                docs.extend(loader.load())
            except Exception:
                pass
        return docs

    def rebuild(self):
        """
        BUILD-TIME STEP — load docs → chunk → embed → store in ChromaDB.

        This is the expensive step (calls the embedding model for every chunk).
        It only needs to run once, or when new documents are added to data/.
        The resulting index is saved to chroma_db/ and loaded cheaply on every
        subsequent startup.
        """
        docs = self._load_documents()
        if not docs:
            # Create an empty DB so the service can still start and answer
            # general StudyHub questions from the system prompt.
            self.vector_db = Chroma(
                persist_directory=DB_DIR,
                embedding_function=self.embeddings,
            )
            return {
                "status": "empty",
                "chunks": 0,
                "message": "No documents found. Add .pdf / .txt / .md files to data/ and call /rebuild.",
            }

        # Split documents into small overlapping chunks.
        # chunk_size=800: roughly one paragraph — focused enough for retrieval.
        # chunk_overlap=150: neighbouring chunks share 150 chars so sentences
        #                    that fall on a boundary are not lost.
        splitter = RecursiveCharacterTextSplitter(chunk_size=800, chunk_overlap=150)
        chunks = splitter.split_documents(docs)

        # Embed every chunk with the local Ollama model and store in ChromaDB.
        # Each chunk becomes a 768-dim (or model-specific) float vector on disk.
        self.vector_db = Chroma.from_documents(
            documents=chunks,
            embedding=self.embeddings,
            persist_directory=DB_DIR,
        )
        return {
            "status": "ok",
            "documents": len(docs),
            "chunks": len(chunks),
            "message": f"Indexed {len(docs)} document(s) into {len(chunks)} chunks.",
        }

    # ── Runtime: answer a question ────────────────────────────────────────────

    def chat(self, question: str, history: Optional[list] = None):
        """
        RUNTIME STEP — embed question → vector search → LLM generation.

        1. Embed the user's question with the same model used at build time.
        2. Search ChromaDB for the 4 most similar chunks (by cosine similarity).
        3. If no chunk is relevant enough, reject the question immediately.
        4. Otherwise pass the chunks as context to the LLM and generate an answer.
        5. Include recent conversation history so follow-up questions work.
        """
        
        context_text = ""
        sources = []

        if self.vector_db is not None:
            try:
                # For follow-up questions ("what about the token?") the question alone is too
                # vague for retrieval. Prepend the last user message to give context.
                retrieval_query = question
                if history:
                    last_user = next((m["content"] for m in reversed(history) if m.get("role") == "user"), None)
                    if last_user:
                        retrieval_query = f"{last_user} {question}"

                # similarity_search_with_relevance_scores returns (Document, score) pairs.
                # Score is 0–1 where 1 = identical meaning. Below the threshold = off-topic.
                results = self.vector_db.similarity_search_with_relevance_scores(retrieval_query, k=4)

                relevant = [(doc, score) for doc, score in results if score >= RELEVANCE_THRESHOLD]

                # Reject only if: nothing relevant found AND it is a substantive question
                # AND it doesn't look like a general question about the platform itself.
                # Short questions and platform-related questions go through to the LLM.
                q_lower = question.lower()
                is_platform_question = any(signal in q_lower for signal in _PLATFORM_SIGNALS)
                if not relevant and results and len(question.split()) > 4 and not is_platform_question:
                    return "I can only answer questions about StudyHub.", []

                if relevant:
                    context_text = "\n\n".join(doc.page_content for doc, _ in relevant)
                    sources = [
                        {
                            "source": os.path.basename(doc.metadata.get("source", "unknown")),
                            "page": doc.metadata.get("page"),
                        }
                        for doc, _ in relevant
                    ]
            except Exception:
                pass

        # Build the prompt as a list of LangChain messages so multi-turn
        # conversation history is preserved correctly.
        system_content = _SYSTEM
        if context_text:
            system_content += f"CONTEXT:\n{context_text}\n"

        messages = [SystemMessage(content=system_content)]

        # Inject the last 6 messages (= 3 user + 3 assistant exchanges) so the
        # LLM understands follow-up questions like "tell me more about that".
        if history:
            for msg in history[-6:]:
                if msg.get("role") == "user":
                    messages.append(HumanMessage(content=msg["content"]))
                elif msg.get("role") == "assistant":
                    messages.append(AIMessage(content=msg["content"]))

        messages.append(HumanMessage(content=question))

        try:
            response = self.llm.invoke(messages)
            return response.content, sources
        except Exception:
            return "Sorry, I couldn't reach the AI model right now. Please try again in a moment.", []
