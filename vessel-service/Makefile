build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/micro-service/shippy-gomicro/vessel-service proto/vessel/vessel.proto
	# dep 工具赞不可用， 直接手动编译
	GOOS=linux GOARCH=amd64 go build
	docker build -d vessel-service .

run:
	docker run 50052:50052 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service
