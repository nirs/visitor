prog := visitor

linux:
	CGO_ENABLED=0 GOOS=$@ go build -o dist/$@/$(prog)

windows:
	GOOS=$@ go build -o dist/$@/$(prog).exe
