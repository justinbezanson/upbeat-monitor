FROM golang:1.23-alpine

WORKDIR /app

# Install git, curl, bash, and libc6-compat (needed for prebuilt binaries on Alpine)
RUN apk add --no-cache git curl bash libc6-compat

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install Turso CLI
RUN curl -sSfL https://get.tur.so/install.sh | bash
# Add Turso to PATH
ENV PATH="/root/.turso:$PATH"

# Copy the current directory contents into the container at /app
COPY . .

# Default command
CMD ["go", "run", "main.go"]