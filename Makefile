install:
	go build -ldflags "-X main.version=`git log --pretty=format:'%h' -n 1`" -o $$GOPATH/bin/tslogs cmd/main.go
	
deps:
	go get -u github.com/kardianos/govendor
	go get -u golang.org/x/tools/cmd/goimports
	go get
	govendor fetch +external +missing
	govendor add +external 

format:
	goimports -w *.go ./cmd/*.go

install_to_tmp: install
	mkdir -p ./tmp && cp -rf $$GOPATH/bin/tslogs ./tmp/tslogs

github_release: install
	TAG=`git describe --exact-match --tags $(git log -n1 --pretty='%h')` TOKEN=$$GITHUB_TOKEN ./staff/github_release.sh

docker_build:
	docker build -t ezotrank/tslogs .

docker_tty:
	docker run --rm -v `pwd`:/go/src/github.com/ezotrank/tslogs -ti ezotrank/tslogs /bin/bash
