FROM golang:alpine AS builder

WORKDIR /app

COPY . .

CMD go run ./develop/dev03/task.go; sleep 20;
