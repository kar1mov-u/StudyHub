import os

from langchain_google_genai import ChatGoogleGenerativeAI, GoogleGenerativeAIEmbeddings
from langchain_chroma import Chroma
from langchain_community.document_loaders import PyPDFLoader, TextLoader, DirectoryLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.prompts import ChatPromptTemplate

DATA_DIR = os.getenv("DATA_DIR", "data")
DB_DIR = os.getenv("DB_DIR", "chroma_db")

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
        self.embeddings = GoogleGenerativeAIEmbeddings(
            model="models/text-embedding-004",
            google_api_key=api_key,
        )
        self.llm = ChatGoogleGenerativeAI(
            model="gemini-2.5-flash-lite",
            google_api_key=api_key,
            temperature=0.2,
        )
        self.vector_db = None
        self._init_db()

    def _init_db(self):
        if os.path.exists(DB_DIR) and any(os.scandir(DB_DIR)):
            self.vector_db = Chroma(
                persist_directory=DB_DIR,
                embedding_function=self.embeddings,
            )
        else:
            self.rebuild()

    def _load_documents(self):
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

        splitter = RecursiveCharacterTextSplitter(chunk_size=800, chunk_overlap=150)
        chunks = splitter.split_documents(docs)
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
                retriever = self.vector_db.as_retriever(search_kwargs={"k": 4})
                relevant_docs = retriever.invoke(question)
                if relevant_docs:
                    context_section = "Relevant context from the knowledge base:\n" + "\n\n".join(
                        d.page_content for d in relevant_docs
                    )
                    sources = [
                        {
                            "source": os.path.basename(d.metadata.get("source", "unknown")),
                            "page": d.metadata.get("page"),
                        }
                        for d in relevant_docs
                    ]
            except Exception:
                pass

        template = "{system}\n\n{context_section}\n\nUser: {question}\n\nAssistant:"
        prompt = ChatPromptTemplate.from_template(template)
        chain = prompt | self.llm
        response = chain.invoke(
            {"system": _SYSTEM, "context_section": context_section, "question": question}
        )
        return response.content, sources
