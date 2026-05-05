import os

from langchain_google_genai import ChatGoogleGenerativeAI, GoogleGenerativeAIEmbeddings
from langchain_chroma import Chroma
from langchain_community.document_loaders import PyPDFLoader, TextLoader, DirectoryLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.prompts import ChatPromptTemplate

DATA_DIR = os.getenv("DATA_DIR", "data")
DB_DIR = os.getenv("DB_DIR", "chroma_db")

# The system prompt tells the LLM who it is and what StudyHub is.
# This text is injected at the top of every request so Gemini always
# has context about the app even when no documents match the question.
_SYSTEM = (
    "You are a helpful assistant for StudyHub, an academic study management platform.\n\n"
    "StudyHub lets students manage academic modules and study resources:\n"
    "- Modules: Courses organised by department. Each has one or more runs.\n"
    "- Module Runs: A module instance in a semester/year, containing weekly sessions.\n"
    "- Weeks: Individual weeks within a run, each holding uploaded resources.\n"
    "- Resources: Files (PDFs, docs) or external links uploaded to a week, stored in S3.\n"
    "- Flashcard Decks: AI-generated flashcards from uploaded PDFs.\n"
    "- Comments: Students comment on weekly resources and upvote/downvote others.\n"
    "- User Profiles: View your own uploads and those of other students.\n"
    "- Academic Terms: Admin-managed semester + year groups for module runs.\n\n"
    "Navigation: Modules -> Module Detail -> Week Detail (resources, comments, flashcards).\n\n"
    "Answer using the provided context when available. If the context does not help, draw on your "
    "general knowledge of StudyHub's features. Be concise and helpful."
)


class RAGChatbot:
    def __init__(self, api_key: str):
        # Embedding model: converts text into a list of numbers (a vector).
        # Two pieces of text that mean similar things will produce vectors
        # that are close together in space — that's how we do "semantic search"
        # later instead of just keyword matching.
        self.embeddings = GoogleGenerativeAIEmbeddings(
            model="models/text-embedding-004",
            google_api_key=api_key,
        )

        # The LLM (large language model): this is the part that actually reads
        # the retrieved context + the user's question and writes a human answer.
        # temperature=0.2 keeps answers factual and consistent (0 = robotic, 1 = creative).
        self.llm = ChatGoogleGenerativeAI(
            model="gemini-2.5-flash-lite",
            google_api_key=api_key,
            temperature=0.2,
        )

        self.vector_db = None
        self._init_db()

    def _init_db(self):
        # On startup: if we already ran /rebuild before, the chroma_db/ folder
        # exists on disk — just open it. Otherwise build it fresh from data/.
        if os.path.exists(DB_DIR) and any(os.scandir(DB_DIR)):
            # ChromaDB is a vector database that lives on the local filesystem.
            # It stores every document chunk alongside its embedding vector so
            # we can search them later by similarity.
            self.vector_db = Chroma(
                persist_directory=DB_DIR,
                embedding_function=self.embeddings,
            )
        else:
            self.rebuild()

    def _load_documents(self):
        # Read every PDF and .txt file from the data/ folder.
        # LangChain's DirectoryLoader walks the folder recursively and hands
        # back a list of Document objects, each with .page_content and .metadata
        # (e.g. {"source": "data/lecture1.pdf", "page": 3}).
        os.makedirs(DATA_DIR, exist_ok=True)
        docs = []
        for loader_cls, glob in [(PyPDFLoader, "**/*.pdf"), (TextLoader, "**/*.txt")]:
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
        docs = self._load_documents()
        if not docs:
            self.vector_db = Chroma(
                persist_directory=DB_DIR,
                embedding_function=self.embeddings,
            )
            return {
                "status": "empty",
                "chunks": 0,
                "message": "No documents found. Drop PDFs or .txt files into data/ and call /rebuild.",
            }

        # Split each document into small overlapping chunks.
        # We can't feed a 50-page PDF directly to the LLM — it's too long and
        # most of it would be irrelevant to any given question.
        # chunk_size=800 chars  → each chunk is roughly one paragraph
        # chunk_overlap=150     → neighbouring chunks share 150 chars so a
        #                         sentence that falls on a boundary isn't lost
        splitter = RecursiveCharacterTextSplitter(chunk_size=800, chunk_overlap=150)
        chunks = splitter.split_documents(docs)

        # Embed every chunk and store it in ChromaDB.
        # This is the slow one-time step: each chunk is sent to Gemini's
        # embedding model, which returns a 768-number vector. Those vectors
        # are saved to chroma_db/ on disk so we never have to redo this
        # unless new documents are added.
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

    def chat(self, question: str):
        context_section = ""
        sources = []

        if self.vector_db is not None:
            try:
                # Retriever: embeds the user's question into a vector,
                # then searches ChromaDB for the 4 chunks whose vectors are
                # most similar (cosine similarity). These are the parts of
                # the documents most likely to contain the answer.
                retriever = self.vector_db.as_retriever(search_kwargs={"k": 4})
                relevant_docs = retriever.invoke(question)

                if relevant_docs:
                    # Paste the raw text of those chunks together.
                    # This becomes the "context" we hand to the LLM so it
                    # answers from real document content, not from memory.
                    context_section = "Relevant context from the knowledge base:\n" + "\n\n".join(
                        d.page_content for d in relevant_docs
                    )
                    # Also record where each chunk came from so the frontend
                    # can show "Sources: lecture1.pdf p.4" below the answer.
                    sources = [
                        {
                            "source": os.path.basename(d.metadata.get("source", "unknown")),
                            "page": d.metadata.get("page"),
                        }
                        for d in relevant_docs
                    ]
            except Exception:
                pass

        # Build the final prompt that gets sent to Gemini:
        #   1. System instructions (who the bot is)
        #   2. Retrieved document chunks (the "retrieved" part of RAG)
        #   3. The user's question
        # Gemini reads all three and writes a grounded answer.
        template = "{system}\n\n{context_section}\n\nUser: {question}\n\nAssistant:"
        prompt = ChatPromptTemplate.from_template(template)

        # The | operator chains the prompt and the LLM together into a pipeline.
        # prompt.invoke(...) formats the template, then passes the result to llm.invoke(...).
        chain = prompt | self.llm
        response = chain.invoke(
            {"system": _SYSTEM, "context_section": context_section, "question": question}
        )
        return response.content, sources
