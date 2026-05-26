# 1. 組み立て用（ビルド環境）の最小のGo言語を用意
FROM golang:1.22-alpine AS builder
WORKDIR /app

# プログラム本体だけをコピー
COPY . .

# 💡 ここがポイント！ビルドする「直前」に、ここで go.mod を自動生成させる
RUN go mod init app
RUN go mod tidy
RUN go build -o main .

# 2. 本番動作用の超軽量環境を用意
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# 💡 ここを追加：画面のファイルも本番環境へ持っていく！
COPY --from=builder /app/index.html .

# 起動コマンド
CMD ["./main"]
