FROM golang:1.25.2-alpine

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache git make

# 作業ディレクトリにファイルをコピー
COPY go.mod go.sum ./

# 依存関係のダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# プロバイダーのビルド
RUN make build

# デフォルトのコマンド
CMD ["sh"]
