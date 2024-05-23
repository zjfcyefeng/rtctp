include rtctp.env
.DEFAULT_GOAL := help

WORKDIR=$(shell pwd)
CURDIRNAME?=$(shell basename ${CURDIR})
CURDATE?=$(shell date '+%Y%m%d')
CURDATETIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LINK_FLAGS='-s -X main.buildTime=${CURDATETIME}'

#============================================================#
# Helper 
#============================================================#
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
## confirm: comfirm operation
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n]' && read ans && [ $${ans:-n} = y ]


#============================================================#
# Git Push
#============================================================#
## gitpush: push code to git server
.PHONY: gitpush
gitpush: 
	@if [ ! $M ]; then \
		echo "add param: M=<your comment for this commit>"; \
		exit 1; \
	fi
	git commit -a -m "${M}"
	git push origin


#============================================================#
# CTP Library
#============================================================#
## ctp/lib: copy ctp lib to lib directory
.PHONY: ctp/lib
ctp/lib:
	mkdir -p /usr/local/rtctp/lib
	cp /root/go/pkg/mod/gitee.com/haifengat/goctp/v2\@v2.0.8/lib/*.so /usr/local/rtctp/lib

#============================================================#
# Build
#============================================================#
## build/cmd: build the cmd application
.PHONY: build/cmd
build/cmd:
	@echo 'Building CMD Application ...'
	go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/cmd/client/rtctp ${WORKDIR}/cmd/client/rtctp/main.go
## build/rest: build the rest api service
.PHONY: build/rest
build/rest:
	@echo 'Building All Restful API Service ...'
	@echo 'Building XXX Restful Service ...'
	go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/rest/xxx ${WORKDIR}/rest/xxx/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/linux_amd64/rest/xxx ${WORKDIR}/rest/xxx/main.go
## build/rpc: build the rpc service
.PHONY: build/rpc
build/rpc:
	@echo 'Building All RPC Service ...'
	@echo 'Building XXX RPC Service ...'
	go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/rpc/xxx ${WORKDIR}/rpc/xxx/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/linux_amd64/rpc/xxx ${WORKDIR}/rpc/xxx/main.go
## build/rtctp: build the real-time ctp application
.PHONY: build/rtctp
build/ctp:
	@echo 'Building Real Time CTP Application ...'
	go build -ldflags=${LINK_FLAGS} -o ${WORKDIR}/bin/rtctp ${WORKDIR}/main.go
