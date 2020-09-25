FROM grpc/go as protoc-builder
WORKDIR /work
COPY Makefile .
COPY interface/interface.proto ./interface/interface.proto
COPY pkg/interface pkg/interface
RUN ["make", "clean", "protobuf"]

FROM golang:1.15.2-alpine3.12 as go-builder
WORKDIR /work
COPY --from=protoc-builder /work/pkg/interface /work/pkg/interface
COPY . .
RUN ["sh", "./scripts/build.sh", "unstacker", "linux", "amd64"]

FROM alpine:3.12.0
COPY --from=go-builder /work/bin /app/bin
CMD ["/app/bin/linux/amd64/unstacker"]