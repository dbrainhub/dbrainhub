#/bin/sh

echo "generate http.pb.go..."
protoc  -I ./ --proto_path=../ --go_out=../ ./http.proto --go_opt=Mgoogle/api/http.proto=google.golang.org/genproto/googleapis/api/annotations --go_opt=Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations
if [ $? -ne 0 ]; then
	echo "generate http.pb.go failed, try install proto and protoc-gen-go."
	echo "proto-gen-go install use: go install github.com/golang/protobuf/protoc-gen-go@v1.5.2"
	exit 1
fi

echo "generate http.swagger.json..."
protoc  -I ./ --proto_path=../ --openapiv2_out=./ ./http.proto --openapiv2_opt=Mgoogle/api/http.proto=google.golang.org/genproto/googleapis/api/annotations --openapiv2_opt=Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations
if [ $? -ne 0 ]; then
	echo "generate http.swagger.json fail, try install proto-gen-openapiv2."
	echo "proto-gen-openapiv2 install use: go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0"
	exit 1
fi

echo "Done."
