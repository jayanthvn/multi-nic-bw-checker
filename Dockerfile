# Build stage
FROM debian:bullseye AS builder
RUN apt-get update && apt-get install -y golang
WORKDIR /app
COPY . .
RUN go build -o bandwidth-checker ./cmd/main.go

# Final stage
FROM debian:bullseye
RUN apt-get update && apt-get install -y iperf3
COPY --from=builder /app/bandwidth-checker /usr/local/bin/bandwidth-checker

# Debugging: List files in /usr/local/bin/
RUN ls -l /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/bandwidth-checker"]

