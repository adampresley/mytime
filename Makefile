DEFAULT_GOAL := build

build:
	go build -o mt

build-mac: 
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o mt && upx ./mt

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mt.exe && upx ./mt.exe

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o mt && upx ./mt
