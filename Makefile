BINARY_NAME = Conductor

all: build

build:
	go get -t -v
	go build -a -o $(BINARY_NAME) -v

run:
	go get -t -v
	go build -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

clean:
	rm $(BINARY_NAME)*