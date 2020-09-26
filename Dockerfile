FROM grpc/go as protoc-builder
WORKDIR /work
COPY Makefile .
COPY interface/interface.proto ./interface/interface.proto
COPY pkg/interface pkg/interface
RUN ["make", "clean", "protobuf"]

FROM golang:1.15.2-alpine3.12 as go-builder
WORKDIR /work
COPY --from=protoc-builder /work/pkg/ .
COPY . .
RUN go mod vendor
RUN sh ./scripts/build.sh unstacker linux amd64

FROM alpine:3.12.0
COPY --from=go-builder /work/bin/linux/amd64 /app/bin
CMD ["/app/bin/unstacker"]