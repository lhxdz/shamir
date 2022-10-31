# 编译选项
# make 默认编译选项
# make clean 清理
# make fmt 文件格式化
NAME=shamir
OUTPUT=bin
COVERFILE=coverage.out
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on

default:
	go mod tidy
	go build -o ./${OUTPUT}/${NAME} ./cmd/shamir.go
	cp -rf ./conf ./${OUTPUT}

clean:
	rm -f ${COVERFILE}
	rm -rf ./${OUTPUT}
	go clean
	go mod tidy

fmt:
	go fmt ./pkg/... ./cmd/...

test:
	go test -v -coverprofile=${COVERFILE} -race ./cmd/... ./pkg/...