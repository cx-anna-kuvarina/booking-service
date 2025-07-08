FROM golang:1.23-alpine AS build
# Copy the source from the current directory to the Working Directory inside the container
WORKDIR /app

#Copy go mod and sum files
COPY go.mod .
COPY go.sum .
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (cache if not changes are made)
RUN --mount=type=cache,target=/root/.cache/go-build  DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/service cmd/*

# Runtime image
FROM alpine:latest

# Copy the binary from the build stage
COPY --from=build /app/bin/service /app/bin/service

WORKDIR /app

ENTRYPOINT ["/app/bin/service"]