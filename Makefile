path = bin/splitter-mac

build:
	cd src && GOOS=darwin GOARCH=amd64 go build -o ../$(path) main.go

install: build
	mv $(path) /usr/local/bin/splitter