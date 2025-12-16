FROM golang:1.22 AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/ssot-devices
RUN go build -o /out/ssot-devices ./cmd/ssot_devices

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/ssot-devices /ssot-devices
EXPOSE 8082
ENTRYPOINT ["/ssot-devices"]
