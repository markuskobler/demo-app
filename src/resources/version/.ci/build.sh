#!/bin/sh
set -eu

export GOPATH=${PWD}/go
export PATH=${GOPATH}/bin:${PATH}

src=${GOPATH}/src/resources/version
out=${PWD}/build

cd ${src}

_test() {
    go test -coverprofile=coverage.out ./...
}

_build() {
    export CGO_ENABLED=0

    echo ">>> Build version-resource"
    go build -o ${out}/version-resource .

    echo "version-resource" > ${out}/tag

    mkdir -p ${out}/etc

    cat <<EOF > ${out}/etc/passwd
root:x:0:0:root:/:/dev/null
nobody:x:65534:65534:nogroup:/:/dev/null
EOF

    cat <<EOF > ${out}/etc/group
root:x:0:
nogroup:x:65534:
EOF

    mkdir -p ${out}/etc/ssl/certs
    cp /etc/ssl/certs/ca-certificates.crt ${out}/etc/ssl/certs/ca-certificates.crt

    cp -r ${src}/Dockerfile ${src}/resource ${out}
}

_test
_build
