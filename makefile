build: 
	go build \
		-ldflags "-X main.buildVersion=$(shell git describe --tags) -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(shell git rev-parse HEAD)'" \
		-o cmd/server/server cmd/server/main.go
	go build \
		-ldflags "-X main.buildVersion=$(shell git describe --tags) -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(shell git rev-parse HEAD)'" \
		-o cmd/agent/agent cmd/agent/main.go
	go build -o cmd/gencert/gencert cmd/gencert/main.go
proto: 
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/proto/octo.proto
test:
	go test -v ./... -short
cover:
	rm -fr coverage
	mkdir coverage
	# флаг -count=1 не использовать кэш
	go test -count=1 ./... -coverprofile coverage/cover.out
	go tool cover -html coverage/cover.out -o coverage/cover.html
