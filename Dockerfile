FROM golang:1.25.4-alpine@sha256:d3f0cf7723f3429e3f9ed846243970b20a2de7bae6a5b66fc5914e228d831bbb AS builder

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

COPY --from=ghcr.io/tarampampam/microcheck:1.2.0@sha256:42cdee55eddc4c5b5bd76cb6df5334ca22ce7dc21066aeea7d059898ffbf84fd /bin/httpcheck /bin/httpcheck

COPY --from=builder /app/dist/template-go /bin/template-go

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT ["/bin/template-go"]
