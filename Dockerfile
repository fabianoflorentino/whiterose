FROM golang:alpine3.22 AS base

WORKDIR /whiterose

COPY . .

RUN apk update -y --no-cache \
  && apk upgrade -y --no-cache \
  && go mod download \
  && GOFLAGS="-trimpath" CGO_DISABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -o /usr/local/bin/whiterose .

FROM base AS development

RUN go install github.com/air-verse/air@latest

COPY --from=base /whiterose/.env.example /root/.env

ENTRYPOINT [ "/go/bin/air" ]
CMD [ "-c", "/whiterose/.air.toml" ]


FROM gcr.io/distroless/static:nonroot AS production

COPY --from=base /usr/local/bin/whiterose /usr/local/bin/whiterose

USER nonroot:nonroot

ENTRYPOINT [ "/usr/local/bin/whiterose" ]
