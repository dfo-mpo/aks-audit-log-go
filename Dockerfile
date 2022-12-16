# syntax=docker/dockerfile:1
#  ╭──────────────────────────────────────────────────────────╮
#  │                       build stage                        │
#  ╰──────────────────────────────────────────────────────────╯
FROM golang:1.19.4-alpine3.17@sha256:f33331e12ca70192c0dbab2d0a74a52e1dd344221507d88aaea605b0219a212f as build
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -u 1001 -D -g '' appuser

WORKDIR /app

COPY . .
# Fetch dependencies.
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /aks-audit-log-go .

#  ╭──────────────────────────────────────────────────────────╮
#  │                       final stage                        │
#  ╰──────────────────────────────────────────────────────────╯
FROM scratch

LABEL org.opencontainers.image.source="https://github.com/JuanPabloSGU/aks-audit-log-go"
LABEL org.opencontainers.image.description="Forward AKS audit logs to falco"
LABEL maintainer="Alexandre.Brassard-Desjardins@dfo-mpo.gc.ca"

# Import from build.
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=build /aks-audit-log-go /aks-audit-log-go

# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENTRYPOINT ["/aks-audit-log-go"]
