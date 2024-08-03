build: 
	go build -o cmd/agent/agent cmd/agent/main.go
	go build -o cmd/server/server cmd/server/main.go
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
