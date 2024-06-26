.PHONY: buildrun
buildrun:
	docker-compose build
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: genMock
genMock:
	mockgen -source=internal/auth/repository.go \
	-destination=internal/auth/mocks/mock_repository.go
	mockgen -source=internal/banner/repository.go \
    	-destination=internal/banner/mocks/mock_repository.go

.PHONY: test
test:
	go test ./test/...


.PHONY: linter
linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
	golangci-lint run -c ./linter/.golangci.yml