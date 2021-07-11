VERSION=`git describe --tags`
BUILD=`date "+%Y%m%d"`

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.build=${BUILD}"

build:
	go build -o bin/hx ${LDFLAGS} ./cmd/hx/hx.go


install:
	go install ${LDFLAGS} ./cmd/hx/hx.go
