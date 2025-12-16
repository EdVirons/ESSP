FROM golang:1.22 AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/ims-api
RUN go build -o /out/ims-api ./cmd/api

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/ims-api /ims-api
EXPOSE 8080
ENTRYPOINT ["/ims-api"]
