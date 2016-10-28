.PHONY: test 

test:
	go test

install:
	go build -o $$GOPATH/bin/tslogs cmd/main.go
	
format:
	go get golang.org/x/tools/cmd/goimports
	goimports -w *.go ./cmd/*.go

install_to_tmp: install
	mkdir -p ./tmp && cp -rf $$GOPATH/bin/tslogs ./tmp/tslogs

docker_build:
	docker build -t ezotrank/tslogs .

binary: docker_build
	mkdir -p ./tmp
	docker create --name tslogs_tmp ezotrank/tslogs
	docker cp tslogs_tmp:/go/bin/tslogs ./tmp/tslogs
	chmod +x ./tmp/tslogs
	docker rm -v tslogs_tmp

docker_tty:
	docker run --rm -v `pwd`:/go/src/github.com/ezotrank/tslogs -ti --entrypoint /bin/bash ezotrank/tslogs 
