FROM golang:1.7.1

COPY . $GOPATH/src/github.com/ezotrank/tslogs

WORKDIR $GOPATH/src/github.com/ezotrank/tslogs
RUN make install

VOLUME /logs
VOLUME /configs

ENTRYPOINT ["/go/bin/tslogs"]
