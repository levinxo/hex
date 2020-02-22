MODULE_NAME = hex
all:
	go build
	./$(MODULE_NAME)
install:
	go install

build:
	go build
	./$(MODULE_NAME)

