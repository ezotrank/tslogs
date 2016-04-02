build: clean init
	GOPATH=`pwd`/vendor go build  -ldflags "-X main.version=`git log --pretty=format:'%h' -n 1`" -o bin/tslogs -race cmd/main.go

deps: init
	(export GOPATH=`pwd`/vendor && cd vendor/src/github.com/ezotrank/tslogs && go get)

init:
	mkdir -p vendor/src/github.com/ezotrank
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
	gzip -5 < bin/tslogs > bin/tslogs.gz
	scp ./bin/tslogs.gz $$USER@$$HOST:"~/tslogs/tslogs.gz"
	ssh $$USER@$$HOST "cd ~/tslogs && gzip -fd tslogs.gz && chmod +x tslogs"

generate_monit_conf:
	@cat staff/monit.conf|sed -e "s/USER/$$USER/g"|sed -e "s/PROJECT/$$PROJECT/g"
