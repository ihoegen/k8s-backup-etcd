IMAGE = ihoegen/backup-etcd
TAG = 1.0.0
test:
	@echo "No tests"
deploy-aws:
	@make build-lambda
	@bash deploy/aws/deploy.sh
build-lambda:
	@GOOS=linux go build -o build/bin/backup-etcd cmd/lambda/*.go
build-docker:
	@GOOS=linux go build -o build/bin/backup-etcd cmd/backup-etcd/*.go
	@docker build -t ${IMAGE}:${TAG} .
build-local:
	go build -o build/bin/backup-etcd cmd/backup-etcd/*.go
push:
	@docker push ${IMAGE}:${TAG}