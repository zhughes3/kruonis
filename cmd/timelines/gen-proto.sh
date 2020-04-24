#!/bin/bash -x

base_filename="timelines"
proto_file=$base_filename.proto

mkdir -p models gen/swagger

# generate stubbed code
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/gogo/protobuf/gogoproto/ \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/ \
  --gogo_out=plugins=grpc:models \
  $proto_file

# generate reverse proxy
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/ \
  --grpc-gateway_out=logtostderr=true:models \
  $proto_file

# generate swagger
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/ \
  --swagger_out=logtostderr=true:./gen/swagger \
  $proto_file

# generate mocks
mockgen -source=models/$base_filename.pb.go -package=models -destination=models/$base_filename.pb.mock.go

