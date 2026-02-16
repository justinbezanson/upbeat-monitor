FROM golang:1.24-bullseye

WORKDIR /app

# Install git, curl, bash, and build tools for CGO
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    bash \
    gcc \
    make \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install Turso CLI
RUN curl -sSfL https://get.tur.so/install.sh | bash
# Add Turso to PATH
ENV PATH="/root/.turso:$PATH"

ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-ldl"

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
