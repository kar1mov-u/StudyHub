"""
BUILD-TIME script — indexes documents into ChromaDB.

Called automatically by the Docker entrypoint on first container start.
re-indexing so startup stays fast.

Can also be run manually to force a full rebuild:
    python3 build_index.py --force
"""

import os
import sys

DB_DIR = os.getenv("DB_DIR", "chroma_db")

force = "--force" in sys.argv

# Skip if the index already exists and --force was not passed.
# This is the key behaviour: build time runs once, runtime just loads.
if not force and os.path.exists(DB_DIR) and any(os.scandir(DB_DIR)):
    print(f"Vector index already exists in {DB_DIR}/ — skipping build. Use --force to rebuild.")
    sys.exit(0)

from rag import RAGChatbot

print("StudyHub RAG — Building vector index...")
bot = RAGChatbot()
result = bot.rebuild()

print(f"Status  : {result['status']}")
print(f"Message : {result['message']}")
if result.get("documents"):
    print(f"Docs    : {result['documents']}")
    print(f"Chunks  : {result['chunks']}")
