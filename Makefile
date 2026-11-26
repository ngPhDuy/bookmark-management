COVERAGE_EXCLUDE=mocks|main.go|test
COVERAGE_THRESHOLD=80

IMAGE ?= bookmark-management
REGISTRY_IMAGE ?= duyngph1/bookmark-service
CONTAINER ?= bookmark-management
TAG ?= latest
REDIS_HOST ?= host.docker.internal
REDIS_PORT ?= 6379
REDIS_PASSWORD ?=

.PHONY: test docs mocks docker-build docker-run docker-stop docker-push

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docs:
	swag init -g cmd/api/main.go -d ./ -o cmd/api/docs

mocks:
	mockery

docker-build:
	docker build -t $(IMAGE):$(TAG) .

docker-run:
	docker run -d --rm --name $(CONTAINER) -p 8080:8080 \
		-e GIN_MODE=release \
		-e REDIS_HOST=$(REDIS_HOST) \
		-e REDIS_PORT=$(REDIS_PORT) \
		-e REDIS_PASSWORD=$(REDIS_PASSWORD) \
		$(IMAGE):$(TAG)

docker-stop:
	docker stop $(CONTAINER)

docker-push: docker-build
	docker tag $(IMAGE):$(TAG) $(REGISTRY_IMAGE):$(TAG)
	docker push $(REGISTRY_IMAGE):$(TAG)
