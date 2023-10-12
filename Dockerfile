# Stage 1: Build
FROM golang:1.21 AS builder

# Install build-essential for cgo
RUN apt-get update && apt-get install -y build-essential

# Create a working directory
WORKDIR /app

# Copy source code
COPY ./cli ./cli
COPY ./server ./server

# Build the tracetest server and CLI
RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o /app/tracetest-server
RUN cd cli && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o /app/tracetest

# Stage 2: Final Image
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the built binaries from the builder stage
COPY --from=builder /app/tracetest-server .
COPY --from=builder /app/tracetest .

EXPOSE 11633/tcp

ENTRYPOINT ["/tracetest-server", "serve"]
