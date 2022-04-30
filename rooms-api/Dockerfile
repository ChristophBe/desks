FROM golang:1.16-bullseye as build

RUN mkdir /src
ADD . /src/
WORKDIR /src/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /src/app app

ENV SERVER_PORT=8080

EXPOSE 8080

USER nonroot:nonroot

CMD ["./app"]
