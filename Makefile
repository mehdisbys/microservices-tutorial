.PHONY: all test

deps:
	export GO111MODULE=on ; go mod download

all:
	make -C ./driver-location linux-binary
	make -C ./gateway linux-binary
	make -C ./zombie-driver linux-binary

test:
	make -C ./driver-location test-dependencies
	make -C ./driver-location test
	make -C ./gateway test
	make -C ./zombie-driver test


clean:
	rm -f ./driver-location/main
	rm -f ./gateway/main
	rm -f ./zombie-driver/main

deploy:
	docker-compose -f docker-compose.yaml up