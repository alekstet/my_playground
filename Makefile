current_dir := $(dir $(mkfile_path))

default: build

clean:
	@echo "Cleaning up previous build..."
	@rm -f  $(current_dir)bin/*

build: clean
	@echo "Building app..."
	@mkdir go_binary_files
	@mkdir go_files
	@go build -o bin/playgo main.go

test:
	@echo "Testing app..."
	@go test ./...
