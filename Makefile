# 编译选项
# make 默认编译选项
# make clean 清理
# make fmt 文件格式化
NAME=shamir
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on

default:
	go mod tidy
	go build -o ./bin/${NAME} ./cmd/shamir.go

clean:
	go clean
	go mod tidy

fmt:
	go fmt ./pkg/... ./cmd/...

test:
	go test -v -coverprofile=coverage.out -race ./cmd/... ./pkg/...