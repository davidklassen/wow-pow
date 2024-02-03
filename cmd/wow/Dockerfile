FROM golang:1.21 AS build

WORKDIR /app

COPY . .

RUN go build -o /wow ./cmd/wow

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /wow /wow

USER nonroot:nonroot

ENTRYPOINT ["/wow"]
