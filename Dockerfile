FROM golang:1.6.2-wheezy

COPY . /go/src/github.com/ezotrank/tslogs
WORKDIR /go/src/github.com/ezotrank/tslogs

RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		curl \
		make \
	&& rm -rf /var/lib/apt/lists/*
RUN make deps && make install

CMD ["tslogs"]