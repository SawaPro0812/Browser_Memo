FROM golang:1.24.1-alpine

WORKDIR /go/src
COPY ./src .

RUN apk upgrade --update && \
    apk --no-cache add git

RUN go get -u github.com/air-verse/air && \
    go build -o /go/bin/air github.com/air-verse/air

CMD ["air", "-c", ".air.toml"]