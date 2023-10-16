 # Stage 1: Build
# FROM golang:1.21 AS builder

# # Set necessary environment variables
# ENV CGO_ENABLED=1
# ENV GOOS=linux
# ENV GOARCH=amd64

# # Install build-essential for cgo
# RUN apt-get update && apt-get install -y build-essential

# # Create a working directory
# WORKDIR /app

# # Copy source code
# COPY ./cli ./cli
# COPY ./server ./server

# # Build the tracetest server and CLI
# RUN cd server && go build -o /app/tracetest-server
# RUN cd cli && go build -o /app/tracetest

# # Stage 2: Final Image
# FROM alpine

# WORKDIR /app

# # Copy the built binaries from the builder stage
# COPY --from=builder /app/tracetest-server /app/tracetest-server
# COPY --from=builder /app/tracetest /app/tracetest

# # Adding /app folder on $PATH to allow users to call tracetest cli on docker
# ENV PATH="$PATH:/app"

# EXPOSE 11633/tcp

# ENTRYPOINT ["/app/tracetest-server", "serve"]


# Stage 1: Build
FROM golang:1.21 AS builder

# Install build-essential for cgo
RUN apt-get update && apt-get install -y build-essential

COPY ./server ./server

RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o /app/tracetest-server
RUN cd cli && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o /app/tracetest

# Stage 2: Final Image
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /app/tracetest-server .
COPY --from=builder /app/tracetest  .

# Remove the proto files from the final stage

RUN rm -rf /app/proto
EXPOSE 11633/tcp
ENTRYPOINT ["/tracetest-server", "serve"]

