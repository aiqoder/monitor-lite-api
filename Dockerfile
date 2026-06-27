# syntax=docker/dockerfile:1

FROM node:20-alpine AS ui-builder
WORKDIR /app/ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ .
ARG VERSION=dev
ENV VITE_APP_VERSION=${VERSION}
RUN npm pkg set "version=${VERSION}" && npm run build-only

FROM golang:1.25-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui-builder /app/ui/dist ./ui/dist
ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X github.com/aiqoder/monitor-lite-api/internal/version.Version=${VERSION}" \
    -o /monitor-lite-api \
    .

FROM alpine:3.21
RUN apk add --no-cache ffmpeg ca-certificates tzdata
WORKDIR /app
COPY --from=go-builder /monitor-lite-api .
EXPOSE 9876
VOLUME ["/app/etc"]
ENTRYPOINT ["./monitor-lite-api"]
CMD ["-f", "etc/tv.yaml"]
