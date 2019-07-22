.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -o ./build/janus-linux-amd64 github.com/qtumproject/janus/cli/janus
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/qtumproject/janus/cli/janus
