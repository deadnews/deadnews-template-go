FROM golang:1.25.0-alpine@sha256:f18a072054848d87a8077455f0ac8a25886f2397f88bfdd222d6fafbb5bba440 AS builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/template-go ./...

FROM gcr.io/distroless/static-debian12@sha256:f2ff10a709b0fd153997059b698ada702e4870745b6077eff03a5f4850ca91b6 AS runtime
LABEL maintainer="deadnews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=builder /app/dist/template-go /bin/template-go

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT ["/bin/template-go"]
