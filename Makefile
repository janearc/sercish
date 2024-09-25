.ONESHELL:

PACKAGE_DIRS=$(shell go list ./... | grep -v /vendor/)
ECR_URI=620055013658.dkr.ecr.us-west-2.amazonaws.com/archeavy/bitey


tidy:
	@go mod tidy

test:
	@echo "crossed fingers emoji running tests"
	@go test -v $(PACKAGE_DIRS)

dockerauth:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 620055013658.dkr.ecr.us-west-2.amazonaws.com

docker:
	docker buildx build --no-cache --tag ${ECR_URI}:latest --platform linux/amd64 --push .

kubelogs:
	kubectl logs `kubectl get pods | grep Running | grep bitey | cut -d ' ' -f 1`

bounce:
	kubectl rollout restart deployment archeavy-bitey-deployment

version:
	@echo "Updating version data"
	@echo "version:" > config/version.yml
	@echo "  build_date: \"`date`\"" >> config/version.yml
	@echo "  build: \"`git describe --tags --always`\"" >> config/version.yml
	@echo "  branch: \"`git branch | grep '^*' | cut -d' ' -f 2`\"" >> config/version.yml
