# syntax=docker/dockerfile:1

FROM golang:1.20.1-alpine
WORKDIR /proxy-checker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./server ./cmd/server && chmod +x ./cmd
RUN echo "proxy-checker started"
CMD ["/proxy-checker/server/"]