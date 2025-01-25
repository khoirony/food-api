FROM golang:1.23.4

# Install dependencies
RUN apt-get update

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Install Go dependencies
RUN go mod tidy

# Run the application
CMD ["go", "run", "main.go"]
