FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/db/migrations ./db/migrations
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]
