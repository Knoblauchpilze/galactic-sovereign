# Default variables
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= server

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build

release: setup
	go build -o build/server main.go

clean:
	rm -rf build
