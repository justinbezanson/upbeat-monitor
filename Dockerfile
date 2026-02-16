FROM golang:1.24-alpine

WORKDIR /app

# Install git, curl, bash, libc6-compat, and build-base (for CGO)
RUN apk add --no-cache git curl bash libc6-compat build-base

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install Turso CLI
RUN curl -sSfL https://get.tur.so/install.sh | bash
# Add Turso to PATH
ENV PATH="/root/.turso:$PATH"

ENV CGO_ENABLED=1

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the current directory contents into the container at /app
COPY . .
COPY . .

# Default command
CMD ["go", "run", "cmd/api/main.go"]
