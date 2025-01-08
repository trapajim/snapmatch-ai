FROM golang:1.23.0-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Install Go tools
RUN go install github.com/air-verse/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest
CMD ["air"]
