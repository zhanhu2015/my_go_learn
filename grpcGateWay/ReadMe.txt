1. 生成hello.pd.gw.go
protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=. \
    hello.proto
2. 生成hello.swagger.json
protoc -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=. \
  hello.proto

3.生成 datafile.go
go-bindata --nocompress -pkg swagger -o datafile.go swag-ui/dist/...