FROM golang:alpine3.22 AS build

WORKDIR /app

COPY . .

RUN GOFLAGS="-trimpath" CGO_DISABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -o /usr/local/bin/whiterose .

FROM gcr.io/distroless/static:nonroot

COPY --from=build /usr/local/bin/whiterose /usr/local/bin/whiterose

USER nonroot:nonroot

ENTRYPOINT [ "/usr/local/bin/whiterose" ]
