FROM golang:1.7

RUN go get github.com/tools/godep
RUN go get github.com/pilu/fresh
