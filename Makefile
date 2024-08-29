BINARY_NAME=hush
MAIN_PACKAGE=./cmd/hush
INSTALL_DIR=/usr/local/bin

build:
	go build -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

clean:
	go clean
	rm -f $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

test:
	go test ./...

.PHONY: build clean install uninstall
