FROM golang:1.7.0

COPY . /go/src/github.com/ezotrank/tslogs
WORKDIR /go/src/github.com/ezotrank/tslogs

RUN make deps && make install

CMD ["tslogs"]