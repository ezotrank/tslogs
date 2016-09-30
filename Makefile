install:
	go build -o $$GOPATH/bin/tslogs cmd/main.go
	
format:
	goimports -w *.go ./cmd/*.go

install_to_tmp: install
	mkdir -p ./tmp && cp -rf $$GOPATH/bin/tslogs ./tmp/tslogs

github_release: install
	TAG=`git describe --exact-match --tags $(git log -n1 --pretty='%h')` TOKEN=$$GITHUB_TOKEN ./staff/github_release.sh

docker_build:
	docker build -t ezotrank/tslogs .

binary: docker_build
	mkdir -p ./tmp
	docker create --name tslogs_tmp ezotrank/tslogs
	docker cp tslogs_tmp:/go/bin/tslogs ./tmp/tslogs
	chmod +x ./tmp/tslogs
	docker rm -v tslogs_tmp

docker_tty:
	docker run --rm -v `pwd`:/go/src/github.com/ezotrank/tslogs -ti ezotrank/tslogs /bin/bash
