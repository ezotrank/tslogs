install:
	go build -ldflags "-X main.version=`git log --pretty=format:'%h' -n 1`" -o $$GOPATH/bin/tslogs cmd/main.go

deps:
	go get -u github.com/kardianos/govendor
	go get -u golang.org/x/tools/cmd/goimports
	go get
	govendor init
	govendor add +external 

github_release:
	TAG=`git describe --exact-match --tags $(git log -n1 --pretty='%h')` TOKEN=$$GITHUB_TOKEN ./staff/github_release.sh

docker_tty:
	docker run --rm -v `pwd`:/go/src/github.com/ezotrank/tslogs -ti ezotrank/tslogs /bin/bash

generate_monit_conf:
	@cat staff/monit.conf|sed -e "s/USER/$$USER/g"|sed -e "s/PROJECT/$$PROJECT/g"
