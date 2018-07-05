all:
	go build

release: darwin linux

darwin: dbp-darwin-amd64

dbp-darwin-amd64: main.go
	GOOS=darwin GOARCH=amd64 go build -o $@ $<

linux: dbp-linux-amd64

dbp-linux-amd64: main.go
	GOOS=linux GOARCH=amd64 go build -o $@ $<

clean:
	rm -rf dbp-*

.PHONY: all release darwin linux clean
