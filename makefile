
VERSION=1.1
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILDTIME=$(shell date "+%Y-%m-%d/%H:%M:%S")
PACKAGE=hz.code/neugls/ads/cmd/ads/ver



# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${PACKAGE}.Version=${VERSION} -X ${PACKAGE}.COMMIT=${COMMIT} -X ${PACKAGE}.BRANCH=${BRANCH} -X ${PACKAGE}.Time=${BUILDTIME}"


build-linux:
	rm dist/* ; \
	GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o dist/hzads-linux-amd64-${VERSION} cmd/ads/main.go ; \
	upx dist/contract-linux-amd64-${VERSION}
build-apple:
	rm dist/* ; \
	GOOS=darwin GOARCH=arm64 go build -v ${LDFLAGS} -o build/hzads-mac-arm64-${VERSION} cmd/ads/main.go 


