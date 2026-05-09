FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source + the knowledge-base documents into the image
COPY main.py rag.py build_index.py ./
COPY data/ ./data/

RUN mkdir -p chroma_db

EXPOSE 8001

# On first start the chroma_db/ volume is empty — build the index then launch.
# On subsequent starts the index already exists so build_index.py is skipped.
CMD ["sh", "-c", "python3 build_index.py && uvicorn main:app --host 0.0.0.0 --port 8001"]
