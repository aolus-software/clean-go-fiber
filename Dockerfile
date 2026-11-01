FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server \
	&& CGO_ENABLED=0 go build -o /out/worker ./cmd/worker \
	&& CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate \
	&& CGO_ENABLED=0 go build -o /out/seed ./cmd/seed

FROM alpine:3.20 AS migrator

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /out/migrate /out/seed ./
USER appuser

CMD ["./migrate", "up"]

FROM alpine:3.20 AS release

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /out/server /out/worker ./
COPY docs/swagger.json ./docs/swagger.json
USER appuser

EXPOSE 3000

CMD ["./server"]
