FROM golang:1.22 AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/sync-worker
RUN go build -o /out/sync-worker ./cmd/worker

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/sync-worker /sync-worker
ENTRYPOINT ["/sync-worker"]
