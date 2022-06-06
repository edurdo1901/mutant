### Testing
test-docker:
	@echo "=> Starting containers and running tests"
	@docker-compose up -d --remove-orphans
	@go test ./... -covermode=atomic -coverpkg=./...
	@docker-compose down -v

test-cover:
	@go test ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./... -count=1
	@go tool cover -html=coverage.out

fmt:	
	@echo "==> Running format"
	go fmt ./...