PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

protobuf: interface/interface.proto | $(PROTOC_GEN_GO)
	protoc --go_out=. interface/interface.proto

compile:
	./scripts/build.sh unstackerd linux/amd64
	./scripts/build.sh unstackerd linux/arm64
	#./scripts/build.sh unstacker linux/amd64

clean:
	echo "Cleaning build directory"
	rm -Rf bin/ || true
	rm pkg/messages/*.pb.go || true

run:
	go run ./cmd/unstackerd -config ./configs/config.yaml

all: clean protobuf compile

.NOTPARALLEL: