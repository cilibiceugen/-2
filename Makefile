.PHONY: build run test clean

BINARY_NAME=depgraph
SRC_DIR=src
TEST_DIR=tests

build:
	cd $(SRC_DIR) && go build -o ../$(BINARY_NAME) .

run:
	cd $(SRC_DIR) && go run . -config ../config.toml

test:
	cd $(SRC_DIR) && go test -v ../$(TEST_DIR)/...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe *.svg
