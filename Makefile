all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64

build-darwin-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/server_darwin_amd64 server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/client_darwin_amd64 client.go

build-darwin-arm64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/server_darwin_arm64 server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/client_darwin_arm64 client.go

build-linux-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/server_linux_amd64 server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/client_linux_amd64 client.go

build-linux-arm64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/server_linux_arm64 server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/client_linux_arm64 client.go

build-windows-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/server_windows_amd64.exe server.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/client_windows_amd64.exe client.go

gox-linux:
	gox -osarch="linux/amd64 linux/arm64" -output="build/server_{{.OS}}_{{.Arch}}"

gox-all:
	gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" -output="build/server_{{.OS}}_{{.Arch}}"

clean:
	rm -f build/server_*