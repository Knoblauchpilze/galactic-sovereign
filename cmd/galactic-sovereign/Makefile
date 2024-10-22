
APPLICATION ?= galactic-sovereign

setup:
	mkdir -p build/bin
	cp configs/galactic-sovereign-template-dev.yml configs/galactic-sovereign-dev.yml

release:
	go build -o build/bin/${APPLICATION} main.go

install: release

run: release
	./build/bin/${APPLICATION} galactic-sovereign-dev

clean:
	rm -rf build configs/galactic-sovereign-dev.yml
