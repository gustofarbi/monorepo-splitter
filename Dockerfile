FROM golang:1.20.3

RUN apt-get -y update && apt-get install -y cmake libssl-dev

RUN cd /tmp && \
    curl -sL https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.tar.gz -o libgit2-1.3.0.tar.gz && \
    tar xzf libgit2-1.3.0.tar.gz && \
    mkdir -p libgit2-1.3.0/build && \
    cd libgit2-1.3.0/build && \
    cmake .. && \
    cmake --build . --target install && \
    cd /tmp && rm -rf libgit2* && \
    ldconfig /usr/local/lib

WORKDIR /go/src/splitter

COPY src/ .

RUN go build -o /usr/local/bin/splitter

WORKDIR /tmp

ENTRYPOINT ["splitter"]
