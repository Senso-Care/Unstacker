
FROM grpc/go as protoc-builder
WORKDIR /work
COPY Makefile .
COPY interface/interface.proto ./interface/interface.proto
COPY pkg pkg
RUN protoc --go_out=. interface/interface.proto

FROM golang:1.15.2-alpine3.12 as go-builder
ARG TARGETPLATFORM
ENV GOPATH="/work"
WORKDIR /work/src/github.com/Senso-Care/Unstacker
COPY --from=protoc-builder /work/pkg/ ./pkg/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN sh ./scripts/build.sh unstacker $TARGETPLATFORM

FROM alpine:3.12.0
ARG TARGETPLATFORM
COPY --from=go-builder /work/src/github.com/Senso-Care/Unstacker/bin/$TARGETPLATFORM /app/bin
CMD ["/app/bin/unstacker"]