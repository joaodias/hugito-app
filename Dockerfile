FROM golang:1.7

RUN wget https://github.com/spf13/hugo/releases/download/v0.17/hugo_0.17-64bit.deb && \
    dpkg -i hugo*.deb

RUN go get github.com/tools/godep
RUN go get github.com/pilu/fresh
