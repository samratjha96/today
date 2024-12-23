# Use a slim Python base image
FROM python:3.11-slim-bullseye AS builder

# Set working directory
WORKDIR /app

# Install system dependencies needed for yfinance and health checks
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install uv
RUN pip install uv

# Copy only the requirements file first for caching
COPY requirements.txt .

# Install project dependencies using uv
RUN uv pip install --system -r requirements.txt

# Copy the application code
COPY . .

# Expose the port your app runs on
EXPOSE 8020

# Set the entrypoint to run your FastAPI app with uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8020"]
