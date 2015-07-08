PROJECT = short
ORGANIZATION = dheckrath
REPO_OWRNER = DanielHeckrath
DOCKER_REPO = $(ORGANIZATION)/$(PROJECT)


SOURCE := $(shell find . -name '*.go')
GOOS := linux
GOARCH := amd64
PROJECT_PATH := github.com.com/$(REPO_OWRNER)/$(PROJECT)

.PHONY=$(PROJECT) docker-build docker-push docker-pull

$(PROJECT): $(SOURCE)
	echo Building for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -it \
	    -v $(shell pwd):/usr/src/go/src/$(PROJECT_PATH) \
		-v $(shell pwd)/Godeps/_workspace/src:/go/src \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/src/go/src/$(PROJECT_PATH) \
	    golang:1.4.2-cross \
	    go build -a -o ./build/$(PROJECT)

build: $(SOURCE)
	godep go build -a -o ./build/$(PROJECT)

run: build
	./build/$(PROJECT)

docker-build: $(PROJECT)
	docker build -t $(DOCKER_REPO) .

docker-push: docker-build
	docker push $(DOCKER_REPO)

docker-pull:
	docker pull $(DOCKER_REPO)

docker-run: docker-build
	docker run \
		-ti \
		--rm \
		-p 8000:8000 \
		-p 8001:8001 \
		-p 8002:8002 \
		--link redis:redis \
		$(DOCKER_REPO)
