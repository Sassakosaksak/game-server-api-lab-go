# Goコンパイラを含む環境で、Linux用のAPI実行ファイルを作る
FROM golang:1.26.5 AS build

WORKDIR /src

# 依存関係の取得をコード変更時の再ビルドから分離する
COPY go.mod ./
RUN go mod download

COPY cmd ./cmd

# 実行用ImageへGoランタイムを持ち込まないため、単体で動くLinux実行ファイルを作る
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# 実行時に必要な最小限のLinux環境だけを使う
FROM alpine:3.22

WORKDIR /app

COPY --from=build /out/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/api"]
