PROJECT = short
ORGANIZATION = dheckrath
REPO_OWRNER = DanielHeckrath
DOCKER_REPO = $(ORGANIZATION)/$(PROJECT)


SOURCE := $(shell find . -name '*.go')
GOOS := linux
GOARCH := amd64
PROJECT_PATH := github.com/$(REPO_OWRNER)/$(PROJECT)

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

build: $(PROJECT)
	docker build -t $(DOCKER_REPO) .

push: docker-build
	docker push $(DOCKER_REPO)

pull:
	docker pull $(DOCKER_REPO)

run: build
	docker run \
		-ti \
		--rm \
		-p 8080:80 \
		-p 8000:8000 \
		-p 8001:8001 \
		--link redis:redis \
		$(DOCKER_REPO)
