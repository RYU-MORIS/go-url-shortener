# 1. 組み立て用（ビルド環境）の最小のGo言語を用意
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 設計図をコピーして中身を組み立てる
COPY go.mod ./
COPY . .
RUN go build -o main .

# 2. 本番動作用の超軽量環境を用意
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

# 起動コマンド
CMD ["./main"]