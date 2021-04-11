FROM golang as builder

RUN wget https://github.com/Kitware/CMake/releases/download/v3.20.1/cmake-3.20.1.tar.gz \
    && tar -xf cmake-3.20.1.tar.gz \
    && cd cmake-3.20.1 \
    && ./bootstrap \
    && make \
    && make install
RUN apt-get update && apt-get install -y libgit2-dev
RUN wget https://github.com/libgit2/libgit2/releases/download/v1.1.0/libgit2-1.1.0.tar.gz \
    && libgit2-1.1.0.tar.gz \
    && cd libgit2-1.1.0

COPY src /go/src
WORKDIR /go/src
#RUN go build -o bin/splitter main.go

#FROM alpine as executor

#COPY --from=builder /tmp/bin/splitter-linux /usr/local/bin

#ENTRYPOINT ["splitter"]