MODULE_NAME = hex
all:
	go build
	./$(MODULE_NAME)
install:
	go install

build:
	go build
	./$(MODULE_NAME)
gomod:
	rm -f go.mod go.sum
	go mod init github.com/levinxo/hex

