# 运行：
```
cd cmd; go run main.go --config ../test/config.json
```

# 问题：
1. goland中显示找不到依赖，但是运行没有问题。
preference-》Go-》Go Module 中勾选 "Enable go modules integration"

2. proto生成：
```
# 安装 protoc-gen-go
brew install protobuf 

# 安装 protoc-gen-go-http
go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
```