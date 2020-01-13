FROM golang:1.13-alpine as build

WORKDIR /src

RUN apk add --no-cache git protobuf && \
    go get github.com/golang/protobuf/protoc-gen-go
#    protoc --proto_path=./ --go_out=plugins=grpc:./ *.proto
