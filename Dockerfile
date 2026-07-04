FROM --platform=${BUILDPLATFORM} golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder

ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOARCH=${TARGETARCH} \
    GOFLAGS="-ldflags=-s" \
    GOCACHE="/cache/build" \
    GOMODCACHE="/cache/mod"

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE} \
    go mod download

COPY --parents cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    go build -o /bin/template-go ./cmd/template-go

FROM gcr.io/distroless/static@sha256:3592aa8171c77482f62bbc4164e6a2d141c6122554ace66e5cc910cadb961ff0 AS runtime

COPY --from=ghcr.io/tarampampam/microcheck:1.4.0@sha256:c9f79cd408626de7c10f2d487d67339f49adf0ba61dde96ede65343269db1f85 /bin/httpcheck /bin/httpcheck

COPY --from=builder /bin/template-go /bin/template-go

USER nonroot:nonroot
HEALTHCHECK NONE
EXPOSE 8000

ENTRYPOINT ["/bin/template-go"]
