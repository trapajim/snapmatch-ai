FROM golang:1.23.0 as builder

WORKDIR /app

# Copy go.mod and go.sum from the root directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY .. .

# Build the Go application
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o cmd/server cmd/main.go

EXPOSE 8080
CMD ["/server"]
# Use a minimal image to run the application
FROM alpine:latest
COPY --from=builder /app/cmd/server /server
EXPOSE 8080
CMD /server
