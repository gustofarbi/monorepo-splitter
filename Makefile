path = bin/splitter-mac

build:
	cd src && GOOS=darwin GOARCH=arm64 PKG_CONFIG_PATH=/usr/local/opt/libgit2/lib/pkgconfig CGO_LDFLAGS="-Wl,-rpath,/usr/local/opt/libgit2/lib"  go build -o ../$(path) main.go

install: build
	mv $(path) /usr/local/bin/splitter

install-libgit2:
	cd /tmp && \
	curl -sL https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.tar.gz -o libgit2-1.3.0.tar.gz && \
	tar xzf libgit2-1.3.0.tar.gz && \
	mkdir -p libgit2-1.3.0/build && \
	cd libgit2-1.3.0/build && \
	cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/opt/libgit2 && \
	sudo cmake --build . --target install && \
	cd /tmp && rm -rf libgit2*

download:
	curl -sL https://github.com/myposter-de/monorepo-splitter/releases/download/0.0.3/splitter-mac -o /usr/local/bin/splitter

test:
	cd src && GOOS=darwin GOARCH=amd64 PKG_CONFIG_PATH=/usr/local/opt/libgit2/lib/pkgconfig CGO_LDFLAGS="-Wl,-rpath,/usr/local/opt/libgit2/lib"  go test ./...

build-docker:
	docker build -t splitter .
