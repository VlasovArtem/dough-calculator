buildImage: ## Build docker image
	@echo "Building docker image"
	@docker build -t bitt/loan -f ./build/Dockerfile .

runTests: ## Run tests
	@echo "Running tests"
	@go test -json ./...

runDockerIntegrationTests: ## Run integration tests
	@echo "Running integration tests"
	@go test -json -tags=integration,docker ./...

runDockerIntegrationTestsReport: ## Run integration tests and generate report
	@echo "Running integration tests and generating report"
	@go test -json -tags=integration,docker ./... | go-test-report

generate: ## Generate
	@echo "Generating"
	@go generate ./...

runDockerIntegrationFullSum: ## Run go test sum
	@echo "Running go test sum"
	@gotestsum -- -tags=integration,docker ./...

generateOpenapiClient: ## Generate openapi client
	@echo "Generating openapi client"
	@oapi-codegen -package=integration_test -o=internal/app/integration_test/app_client_test.go ./api/openapi.yaml

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
