FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o bandwidth-checker ./cmd/main.go

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y iperf3 && apt-get clean
COPY --from=builder /app/bandwidth-checker /usr/local/bin/bandwidth-checker
ENTRYPOINT ["/usr/local/bin/bandwidth-checker"]

