
APPLICATION ?= stellar-dominion

setup:
	mkdir -p build/bin
	cp configs/stellar-dominion-template-dev.yml configs/stellar-dominion-dev.yml

release:
	go build -o build/bin/${APPLICATION} main.go

install: release

run: release
	./build/bin/${APPLICATION} stellar-dominion-dev

clean:
	rm -rf build configs/stellar-dominion-dev.yml
