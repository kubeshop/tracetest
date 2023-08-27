# Stage 1: Build
FROM golang:1.21 AS builder

# Set necessary environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Install build-essential for cgo
RUN apt-get update && apt-get install -y build-essential

# Create a working directory
WORKDIR /app

# Copy source code
COPY ./cli ./cli
COPY ./server ./server

# Build the tracetest server and CLI
RUN cd server && go build -o /app/tracetest-server
RUN cd cli && go build -o /app/tracetest

# Stage 2: Final Image
FROM alpine

WORKDIR /app

# Copy the built binaries from the builder stage
COPY --from=builder /app/tracetest-server /app/tracetest-server
COPY --from=builder /app/tracetest /app/tracetest

# Adding /app folder on $PATH to allow users to call tracetest cli on docker
ENV PATH="$PATH:/app"

EXPOSE 11633/tcp

ENTRYPOINT ["/app/tracetest-server", "serve"]
