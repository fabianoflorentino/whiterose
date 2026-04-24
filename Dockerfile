FROM golang:1.26-alpine3.23 AS base

WORKDIR /whiterose

COPY . .

RUN apk add --no-cache git && \
    go mod download && \
    GOFLAGS="-trimpath" CGO_DISABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -o /usr/local/bin/whiterose .

FROM base AS development

RUN go install github.com/air-verse/air@latest && \
    rm -rf /go/pkg/mod

COPY --from=base /whiterose/.env.example /root/.env

ENTRYPOINT [ "/go/bin/air" ]
CMD [ "-c", "/whiterose/.air.toml" ]


FROM gcr.io/distroless/static:nonroot AS production

COPY --from=base /usr/local/bin/whiterose /usr/local/bin/whiterose

USER nonroot:nonroot

ENTRYPOINT [ "/usr/local/bin/whiterose" ]
