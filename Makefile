build: clean init
	GOPATH=`pwd`/vendor go build -o build/tslogs -race bin/main.go

deps: init
	(export GOPATH=`pwd`/vendor && cd vendor/src/github.com/ezotrank/tslogs && go get)

init:
	mkdir -p vendor/ssdlrc/github.com/ezotrank
	rm -rf vendor/src/github.com/ezotrank/tslogs
	ln -snf `pwd` vendor/src/github.com/ezotrank/tslogs

clean:
	mkdir -p build
	rm -rf build/*

clean_tmp:
	rm -rf ./tmp
	mkdir ./tmp

format: init
	gofmt -w *.go
	goimports -w *.go

deploy: build
	ssh $$USER@$$HOST "mkdir -p ~/tslogs"
	gzip -5 build/tslogs
	scp ./build/tslogs.gz $$USER@$$HOST:"~/tslogs/tslogs.gz"
	ssh $$USER@$$HOST "cd ~/tslogs && gzip -fd tslogs.gz"
