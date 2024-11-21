FROM nekkkkitch/docker AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -a -o ses ./services/ses/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ses .
CMD ["./ses"]