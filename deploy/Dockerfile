# Клієнтський застосунок: webpack 5 цього проєкту збирається лише на Node 16
# (на Node 18 падає з помилкою OpenSSL). Debian, а не alpine, бо image-webpack-loader
# завантажує готові бінарники mozjpeg / optipng / pngquant, зібрані під glibc.
FROM node:16-bullseye AS web
WORKDIR /web
COPY web/app/package.json web/app/package-lock.json ./
RUN npm ci --legacy-peer-deps --no-audit --no-fund
COPY web/app/ ./
RUN npm run build

# Сервер.
FROM golang:1.17-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/generator ./cmd

# Образ для запуску.
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -u 1000 app
WORKDIR /app
COPY --from=server /out/generator ./generator
COPY --from=web /web/dist ./web/app/dist
COPY web/public-pages ./web/public-pages
# ethers для мінт-сайту віддає сам сервер (/static/js/vendor/), а не сторонній CDN.
COPY --from=web /web/node_modules/ethers/dist/ethers.umd.min.js ./web/public-pages/vendor/ethers.umd.min.js
COPY deploy/config.json ./configs/config.json
RUN mkdir -p data/files && chown -R app:app data
USER app

EXPOSE 9012
VOLUME ["/app/data"]
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD wget -qO /dev/null http://127.0.0.1:9012/ || exit 1

ENTRYPOINT ["./generator", "run", "--config=./configs"]
