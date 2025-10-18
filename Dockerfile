FROM golang:1.25.3-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/template-go ./...

FROM gcr.io/distroless/static-debian12@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c AS runtime
LABEL maintainer="deadnews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=builder /app/dist/template-go /bin/template-go

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT ["/bin/template-go"]
