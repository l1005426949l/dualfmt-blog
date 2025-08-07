.PHONY: proto
proto:
	./bin/protoc --proto_path=. \
		--proto_path=./include \
		--proto_path=./third_party/googleapis \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		--go-errors_out=paths=source_relative:. \
		api/file/v1/file.proto \
		internal/conf/conf.proto

.PHONY: wire
wire:
	wire
