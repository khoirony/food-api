FROM golang:1.23.4

# Install dependencies
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libtesseract-dev \
    libleptonica-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Install Go dependencies
RUN go mod tidy

# Run the application
CMD ["go", "run", "main.go"]
