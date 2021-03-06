REPO         = max-grape/test-backend
CWD          = /go/src/github.com/$(REPO)
IMAGE        = ghcr.io/$(REPO):latest
IMAGE_GO     = golang:1.17.8
IMAGE_ALPINE = alpine:3.15.2

lint:
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) -e GOFLAGS=-mod=vendor \
		golangci/golangci-lint:v1.45 golangci-lint run

unit:
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) -e GOFLAGS=-mod=vendor \
		$(IMAGE_GO) sh -c "go list ./... | grep -v 'vendor\|test' | xargs go test -race -v"

build:
	@docker build \
		--build-arg IMAGE_GO=$(IMAGE_GO) \
		--build-arg IMAGE_ALPINE=$(IMAGE_ALPINE) \
		--build-arg CWD=$(CWD) \
		--build-arg GOOS=linux \
		--build-arg GOARCH=amd64 \
		-t $(IMAGE) .

acceptance: down
	@IMAGE=$(IMAGE) IMAGE_GO=$(IMAGE_GO) CWD=$(CWD) docker-compose -f ./test/docker-compose.acceptance.yml up -d --scale acceptance=0
	@IMAGE=$(IMAGE) IMAGE_GO=$(IMAGE_GO) CWD=$(CWD) docker-compose -f ./test/docker-compose.acceptance.yml up --abort-on-container-exit acceptance

down:
	@IMAGE=$(IMAGE) docker-compose -f ./test/docker-compose.acceptance.yml down -v --remove-orphans
