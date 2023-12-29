PROJECT:=go-protector

.PHONY: build-current

build-current:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -o go-protector .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o smart_iam_linux .