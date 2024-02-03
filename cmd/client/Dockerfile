FROM golang:1.21 AS build

WORKDIR /app

COPY . .

RUN go build -o /client ./cmd/client

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /client /client

USER nonroot:nonroot

ENTRYPOINT ["/client"]
