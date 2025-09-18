# Build stage
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bizbuy main.go

# Final stage
FROM scratch
WORKDIR /
COPY --from=builder /app/bizbuy /bizbuy
COPY --from=builder /app/templates /templates
EXPOSE 8080
ENTRYPOINT ["/bizbuy"]
