FROM golang:latest AS builder

WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -trimpath -ldflags=-buildid= -o main ./

FROM ghcr.io/greboid/dockerbase/nonroot:1.20250326.0

COPY --from=builder /app/main /dockercleanup
CMD ["/dockercleanup"]
