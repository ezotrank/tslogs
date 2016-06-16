install:
	go build -ldflags "-X main.version=`git log --pretty=format:'%h' -n 1`" -o $$GOPATH/bin/tslogs cmd/main.go

deps:
	go get -u github.com/kardianos/govendor
	go get -u golang.org/x/tools/cmd/goimports
	govendor init
	go get
	govendor add +external 

github_release:
	TAG=$$TAG TOKEN=$$TOKEN ./staff/github_release.sh

docker_tty:
	docker run --rm -v `pwd`:/go/src/github.com/ezotrank/tslogs -ti ezotrank/tslogs /bin/bash

generate_monit_conf:
	@cat staff/monit.conf|sed -e "s/USER/$$USER/g"|sed -e "s/PROJECT/$$PROJECT/g"
