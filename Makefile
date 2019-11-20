NUMBERS = 0 1 2 3
.PHONY: build
.PHONY: clean
build:
	for i in $(NUMBERS); do \
		mkdir -p build/node$$i/data; \
		cp build/priv_validator_state.json build/node$$i/data/priv_validator_state.json; \
	done
	docker build . -t webchat:latest

clean:
	for i in $(NUMBERS); do \
		rm -rf build/node$$i/data; \
	done

all: clean build
