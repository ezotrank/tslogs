FROM golang:1.7.1

COPY . $GOPATH/src/github.com/ezotrank/tslogs

WORKDIR $GOPATH/src/github.com/ezotrank/tslogs
RUN make install

VOLUME /logs
VOLUME /configs
ENV logging info

CMD ["/go/bin/tslogs", "-config", "/configs/config.conf",  "-logging", "$logging"]