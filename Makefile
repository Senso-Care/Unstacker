PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

protobuf: interface/interface.proto | $(PROTOC_GEN_GO)
	protoc --go_out=paths=source_relative:./pkg interface/interface.proto

compile:
	./scripts/build.sh unstacker linux/amd64
	./scripts/build.sh unstacker linux/arm64
	#./scripts/build.sh sender linux/amd64

clean:
	echo "Cleaning build directory"
	rm -Rf bin/ || true
	rm pkg/interface/*.pb.go || true

run:
	go run ./cmd/unstacker -config ./configs/config.yaml

all: clean protobuf compile

.NOTPARALLEL: