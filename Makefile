.PHONY: docker
docker:
	@rm webook ||true
	@$ENV:GOOS="linux"
	@$ENV:GOARCH="amd64"
	@go build -o webook .
	@docker rmi -f yang/webook:v0.0.1
	@docker build -t yang/webook:v0.0.1 .