echo "clean old files..."
rm ./http.ph.go
rm ./http.pb.go

set -e

echo "generate http.ph.go..."
protoc  --proto_path=../  -I ./ http.proto --go_opt=Mgoogle/api/http.proto=google.golang.org/genproto/googleapis/api/annotations --go_opt=Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations  --go_out=../
mv http.pb.go http.ph.go

echo "generate http.pb.go..."
protoc  --proto_path=../  -I ./ http.proto --go-http_opt=Mgoogle/api/http.proto=google.golang.org/genproto/googleapis/api/annotations --go-http_opt=Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations  --go-http_out=../
mv http_http.pb.go http.pb.go

echo "success."