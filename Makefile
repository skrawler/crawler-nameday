.PHONY: data

full: fetch-and-generate data

fetch-and-generate:
	go run cmd/download-test-data/main.go
	go run cmd/parse-latest/main.go

data:
	go-bindata -nocompress -nometadata -pkg nameday -o bindata.go data/...

update-deps:
	rm -rf vendor
	dep ensure
	dep ensure -update

test:
	go test -v ./...
