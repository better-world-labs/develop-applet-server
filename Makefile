build:
	make gone
	go build -ldflags="-w -s" -tags musl -o bin/server cmd/server/main.go

run:
	make gone
	go run cmd/server/main.go

gone:
	make install-gone
	make install-mockgen
	go generate ./...

install-gone:
	go install github.com/gone-io/gone/tools/gone@latest

install-mockgen:
	go install github.com/golang/mock/mockgen@latest

gen:
	make gone

watch:
	gone -s internal -p internal -f Priest -o internal/priest.go -w

install-tools:
	make install-gone
	make install-mockgen

build-docker:
	docker build -t moyu .

test:
	go test ./... -v