prog := visitor

all: linux windows

linux:
	CGO_ENABLED=0 GOOS=$@ go build -o dist/$(prog)-$@-amd64

windows:
	GOOS=$@ go build -o dist/$(prog)-$@-amd64.exe
