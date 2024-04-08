
APPLICATION ?= users

setup:
	mkdir -p build build/configs build/bin
	cp configs/users-template-dev.yml configs/users-dev.yml

config:
	cp configs/users-dev.yml build/configs/users-dev.yml

release:
	go build -o build/bin/${APPLICATION} main.go

install: release config

run: release
	./build/bin/${APPLICATION} users-dev

clean:
	rm -rf build configs/users-dev.yml
