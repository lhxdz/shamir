# 编译选项
# make 默认编译选项
# make clean 清理
# make fmt 文件格式化
NAME=shamir
OUTPUT=bin
COVERFILE=coverage.out
CONFIG_DIR=/var/lib/shamir
INSTALL_DIR=/usr/local/bin
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on

default:
	@echo "build binary file"
	go mod tidy
	go build -o ./${OUTPUT}/${NAME} ./cmd/shamir.go
	cp -rf ./conf ./${OUTPUT}

install:default _install clean

_install:
	@echo "installing shamir..."
	cp ./${OUTPUT}/${NAME} ${INSTALL_DIR}/${NAME}
	chmod 755 ${INSTALL_DIR}/${NAME}
	mkdir -p ${CONFIG_DIR}
	cp -f ./${OUTPUT}/conf/* ${CONFIG_DIR}/.

uninstall:
	@echo "uninstalling shamir..."
	rm -f ${INSTALL_DIR}/${NAME}
	rm -rf ${CONFIG_DIR}

clean:
	@echo "clean..."
	go mod tidy
	go clean
	rm -f ${COVERFILE}
	rm -rf ./${OUTPUT}

fmt:
	go fmt ./pkg/... ./cmd/...

test:
	go test -v -coverprofile=${COVERFILE} -race ./cmd/... ./pkg/...