PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

protobuf: interface/interface.proto | $(PROTOC_GEN_GO)
	protoc --go_out=paths=source_relative:./pkg interface/interface.proto

compile:
	./scripts/build.sh mqtt-transformer linux amd64
	./scripts/build.sh mqtt-transformer linux arm64
	./scripts/build.sh sender linux amd64

clean:
	echo "Cleaning build directory"
	rm -Rf bin/ || true
	rm pkg/interface/*.pb.go || true

run:
	go run ./cmd/mqtt-transformer -config ./configs/config.yaml

all: clean protobuf compile

.NOTPARALLEL: