all: say_version build test clean

tests: build test clean

say_version:
	@echo "gominima/docgen v1"

build:
	@echo "Building..."
	go build docgen.go

test:
	@echo "Testing..."
	./docgen tests/test.go

clean:
	@echo "Cleaning up..."
	rm -rf docgen