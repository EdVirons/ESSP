FROM golang:1.22 AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/ssot-parts
RUN go build -o /out/ssot-parts ./cmd/ssot_parts

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/ssot-parts /ssot-parts
EXPOSE 8083
ENTRYPOINT ["/ssot-parts"]
