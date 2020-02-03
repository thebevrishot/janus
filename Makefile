.PHONY: release
release: darwin linux

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/qtumproject/janus/cli/janus

.PHONY: linux
linux:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/qtumproject/janus/cli/janus


	
