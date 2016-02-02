build: clean
	GO15VENDOREXPERIMENT=1 go build -o build/tslogs -race main.go

deps: init
	GO15VENDOREXPERIMENT=1 go get golang.org/x/tools/cmd/goimports

init:
	export GOPATH=`pwd`/vendor
	export PATH=`pwd`/vendor/bin:$PATH
	mkdir -p vendor/src/github.com/ezotrank/tslogs
	ln -snf `pwd`/tslogs vendor/src/github.com/ezotrank/tslogs/tslogs

clean:
	mkdir -p build
	rm -rf build/*

clean_tmp:
	rm -rf ./tmp
	mkdir ./tmp

format: init
	goimports -w ./..
	gofmt -w ./..

package_as: init build
	rm -rf package/tslogs
	mkdir -p package/tslogs
	cp -rf build/tslogs package/tslogs
	cp -rf staff/as_yasen_config.toml package/tslogs/config.toml
	cp -rf staff/as_monit.conf package/tslogs
	(cd package && tar -cvzf tslogs.tar.gz ./tslogs/* && rm -rf ./tslogs)
