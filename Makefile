# 编译选项
# make 默认编译选项
# make clean 清理
# make fmt 文件格式化
NAME=shamir
COVERFILE=coverage.out
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on

default:
	go mod tidy
	go build -o ./bin/${NAME} ./cmd/shamir.go

clean:
	rm -f ${COVERFILE}
	rm -rf ./bin
	go clean
	go mod tidy

fmt:
	go fmt ./pkg/... ./cmd/...

test:
	go test -v -coverprofile=${COVERFILE} -race ./cmd/... ./pkg/...