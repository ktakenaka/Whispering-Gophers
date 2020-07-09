FROM golang

RUN apt-get update && apt-get upgrade -y
RUN go get github.com/campoy/whispering-gophers/hello

WORKDIR /go/src
