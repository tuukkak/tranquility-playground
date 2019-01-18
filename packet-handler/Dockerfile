FROM golang:1.11.4-alpine

RUN apk add --no-cache git openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY ./src /go/src
WORKDIR /go/src
RUN mkdir bin

ENV GOBIN $GOPATH/bin

RUN go get ./

EXPOSE 3000

CMD ["go", "run", "main.go"]