FROM golang:1.21 AS build

WORKDIR /app

COPY . .

RUN go build ./cmd/client

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /app/client /client

USER nonroot:nonroot

ENTRYPOINT ["/client"]
