GO_BUILD_FILE:=./cmd/app/main.go
GO_VER:=1.23
ALPINE_VER:=3.18
VERSION:=1.0.0
DOCKER_IMG="go-clean-template:dev"
ENVIRONMENT="LOCAL"

.PHONY: run lint docker-build docker-run

build:
	go build -o /tmp/app ${GO_BUILD_FILE}

run:
	go run ${GO_BUILD_FILE}

lint:
	golangci-lint run -v --color=always $GO_PACKAGES --timeout 4m

docker-build:
	docker build \
		--build-arg=GO_VER="${GO_VER}" \
		--build-arg=ALPINE_VER="${ALPINE_VER}" \
		--build-arg=VERSION="${VERSION}" \
		-t ${DOCKER_IMG} \
		-f ./docker/Dockerfile .

docker-run: docker-build
	docker run \
		-e environment=${ENVIRONMENT} \
		${DOCKER_IMG}


.PHONY: cover cover100

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

cover100:
	go test -short -count=100 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: mocks

MOCK_SRC_DIRS:=internal/service pkg/logger/ pkg/monitoring
	
MOCK_PREFIX := mocks

mocks:
	@for dir in $(MOCK_SRC_DIRS); do \
		for file in $$dir/*.go; do \
			if echo $$file | grep -q 'test'; then \
                continue; \
            fi; \
			file_name=$$(basename $$file); \
			mock_name=$$(echo $$file_name | sed 's/\.go/_mock.go/'); \
			mockgen -destination=$$dir/$(MOCK_PREFIX)/$$mock_name -source=$$file; \
		done \
	done
	@echo "Mocks generated successfully :3"\
