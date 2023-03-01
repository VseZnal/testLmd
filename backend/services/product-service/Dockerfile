FROM golang:1.20-alpine AS builder

WORKDIR /testLmd

COPY ./backend/ ./backend

RUN cd backend && go mod download
RUN cd ./backend/services/product-service/ && go build -o /compiled

FROM alpine:3.16 AS runtime

WORKDIR /app

COPY --from=builder /compiled /app

ENTRYPOINT ["/app/compiled"]
