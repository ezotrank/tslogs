FROM golang:1.7.1

COPY . /go/src/github.com/ezotrank/tslogs
WORKDIR /go/src/github.com/ezotrank/tslogs

RUN make install
VOLUME /logs
VOLUME /configs

CMD ["/go/bin/tslogs"]