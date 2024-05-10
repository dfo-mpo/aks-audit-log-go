# syntax=docker/dockerfile:1
#  ╭──────────────────────────────────────────────────────────╮
#  │                       build stage                        │
#  ╰──────────────────────────────────────────────────────────╯

FROM golang:1.22.3-alpine3.19@sha256:2a882244fb51835ebbd8313bffee83775b0c076aaf56b497b43d8a4c72db65e1 as build

# Install necessary packages
RUN apk update && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates && \
    apk add --no-cache shadow

# Create appuser
RUN groupadd -g 1001 appuser \
    && useradd -u 1001 -g 1001 -m appuser

WORKDIR /app

COPY . .
# Fetch dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /aks-audit-log-go .
#  ╭──────────────────────────────────────────────────────────╮
#  │                       final stage                        │
#  ╰──────────────────────────────────────────────────────────╯
FROM scratch

LABEL org.opencontainers.image.source="https://github.com/dfo-mpo/aks-audit-log-go"
LABEL org.opencontainers.image.description="Forward AKS audit logs to falco"
LABEL maintainer="Alexandre.Brassard-Desjardins@dfo-mpo.gc.ca"

# Import from build
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=build /aks-audit-log-go /aks-audit-log-go

# Use an unprivileged user
USER appuser

# Expose port 9000
EXPOSE 9000

ENTRYPOINT ["/aks-audit-log-go"]
