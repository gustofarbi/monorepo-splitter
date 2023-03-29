build:
	GOOS=darwin GOARCH=amd64 PKG_CONFIG_PATH=/usr/local/opt/libgit2/lib/pkgconfig CGO_LDFLAGS="-Wl,-rpath,/usr/local/opt/libgit2/lib" \
	go build main.go

test:
	GOOS=darwin GOARCH=amd64 PKG_CONFIG_PATH=/usr/local/opt/libgit2/lib/pkgconfig CGO_LDFLAGS="-Wl,-rpath,/usr/local/opt/libgit2/lib" \
	go test ./...

run:
	GOOS=darwin GOARCH=amd64 PKG_CONFIG_PATH=/usr/local/opt/libgit2/lib/pkgconfig CGO_LDFLAGS="-Wl,-rpath,/usr/local/opt/libgit2/lib" \
	go run main.go -c splitter.yaml
