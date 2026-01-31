FROM golang:1.25.6-alpine@sha256:98e6cffc31ccc44c7c15d83df1d69891efee8115a5bb7ede2bf30a38af3e3c92 AS builder

ENV GOCACHE="/cache/go-build" \
    CGO_ENABLED=0

WORKDIR /app

COPY --parents go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/template-go /app/cmd/template-go

FROM gcr.io/distroless/static-debian13@sha256:972618ca78034aaddc55864342014a96b85108c607372f7cbd0dbd1361f1d841 AS runtime
LABEL maintainer="deadnews <deadnewsgit@gmail.com>"

COPY --from=ghcr.io/tarampampam/microcheck:1.3.0@sha256:79c187c05bfa67518078bf4db117771942fa8fe107dc79a905861c75ddf28dfa /bin/httpcheck /bin/httpcheck

COPY --from=builder /app/dist/template-go /bin/template-go

USER nonroot:nonroot
HEALTHCHECK NONE
EXPOSE 8000

ENTRYPOINT ["/bin/template-go"]
