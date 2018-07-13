VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.Version=$(VERSION)"'

release: gh-release clean dist
	govendor sync
	github-release release \
	--security-token $$GH_LOGIN \
	--user segmentio \
	--repo chamber-s3 \
	--tag $(VERSION) \
	--name $(VERSION)

	github-release upload \
	--security-token $$GH_LOGIN \
	--user segmentio \
	--repo chamber-s3 \
	--tag $(VERSION) \
	--name chamber-s3-$(VERSION)-darwin-amd64 \
	--file dist/chamber-s3-$(VERSION)-darwin-amd64

	github-release upload \
	--security-token $$GH_LOGIN \
	--user segmentio \
	--repo chamber-s3 \
	--tag $(VERSION) \
	--name chamber-s3-$(VERSION)-linux-amd64 \
	--file dist/chamber-s3-$(VERSION)-linux-amd64

clean:
	rm -rf ./dist

dist:
	mkdir dist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/chamber-s3-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/chamber-s3-$(VERSION)-linux-amd64

gh-release:
	go get -u github.com/aktau/github-release

govendor:
	go get -u github.com/kardianos/govendor
