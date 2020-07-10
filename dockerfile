FROM golang:1.13.0-alpine3.10

RUN mkdir /app

WORKDIR /app

COPY bin/. .

CMD ["/app/server"]
