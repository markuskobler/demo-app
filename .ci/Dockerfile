from golang:1.10-alpine

RUN \
    apk --no-cache add git && \
    go get -u github.com/a-urth/go-bindata/... && \
    mv ./bin/go-bindata /usr/local/bin && \
    rm -rf *
