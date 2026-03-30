FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o weblise-server ./server

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates sqlite-libs

COPY --from=builder /build/weblise-server /app/weblise-server

RUN mkdir -p /app/data

ENV HTTP_PORT=8080
ENV WS_PORT=8443
ENV DB_PATH=/app/data/weblise.db

EXPOSE 8080 8443

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O- http://localhost:8080/health || exit 1

CMD ["./weblise-server"]
