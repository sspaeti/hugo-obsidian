.PHONY: build install clean test

# Build the binary with development name
build:
	go build -o hugo-obsidian-dev

# Install the binary to ~/.go/bin
install: build
	cp hugo-obsidian-dev $(HOME)/.go/bin/hugo-obsidian-dev

install-prod: build
	cp hugo-obsidian-dev $(HOME)/.go/bin/hugo-obsidian

# Run the program with default settings
run: build
	./hugo-obsidian-dev -input=content -output=data -index=true

# Clean up build artifacts
clean:
	rm -f hugo-obsidian-dev

# Run tests
test:
	go test ./...
