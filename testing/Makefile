.PHONY: build build-openzeppelin openzeppelin
build: build-openzeppelin

# Build openzeppelin docker files with docker-compose
openzeppelin-docker-compose:
	./run-openzeppelin-tests.sh

build-openzeppelin:
	docker build -t qtum/janus-openzeppelin -f ./openzeppelin.Dockerfile .

init:
	git submodule init
	cd openzeppelin-contracts && yarn install

openzeppelin:
	cd openzeppelin-contracts && truffle test --network testing

openzeppelin-docker: build
	docker run --rm qtum/janus-openzeppelin