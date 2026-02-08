FROM golang:1.23-alpine

WORKDIR /app

# Install git and curl for development convenience
RUN apk add --no-cache git curl

# Copy the current directory contents into the container at /app
COPY . .

# Default command
CMD ["go", "run", "main.go"]