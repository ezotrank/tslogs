clean:
	mkdir -p build
	rm -rf build/*

build: clean
	go build -o build/tslogs -race main.go

package_as: build
	rm -rf package/tslogs
	mkdir -p package/tslogs
	cp -rf build/tslogs package/tslogs
	cp -rf staff/as_yasen_config.toml package/tslogs/config.toml
	cp -rf staff/as_monit.conf package/tslogs
	(cd package && tar -cvzf tslogs.tar.gz ./tslogs/* && rm -rf ./tslogs)
