FROM golang:1.21 AS build

WORKDIR /app

COPY . .

RUN go build ./cmd/wow

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /app/wow /app/db.txt ./

USER nonroot:nonroot

ENTRYPOINT ["/wow"]
