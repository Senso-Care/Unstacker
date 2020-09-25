FROM grpc/go as protoc-builder
WORKDIR /work
COPY Makefile .
COPY interface/interface.proto ./interface/interface.proto
COPY pkg/interface pkg/interface
RUN ["make", "clean", "protobuf"]

FROM golang:latest as go-builder
WORKDIR /work
COPY --from=protoc-builder /work/pkg/interface /work/pkg/interface
COPY . .
RUN ["make", "compile"]

FROM debian:buster-slim
COPY --from=go-builder /work/bin /app/bin
CMD ["/app/bin/linux/amd64/unstacker"]