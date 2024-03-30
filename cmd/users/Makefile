
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= users

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build
	mkdir -p build/bin

release: setup
	go build -o build/bin/${APPLICATION} main.go

run: release
	./build/bin/${APPLICATION}

clean:
	rm -rf build
