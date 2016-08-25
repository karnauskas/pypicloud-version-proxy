.PHONY: linux

clean:
	rm -f pypicloud-version-proxy

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w"