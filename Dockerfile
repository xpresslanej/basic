FROM golang:alpine

LABEL maintainer="Jeff Lane <xpresslanej@gmail.com>"

WORKDIR $GOPATH/src/github.com/xpresslanej/basic

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["basic"]