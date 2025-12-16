FROM golang:1.22 AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/ssot-school
RUN go build -o /out/ssot-school ./cmd/ssot_school

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/ssot-school /ssot-school
EXPOSE 8081
ENTRYPOINT ["/ssot-school"]
