.PHONY: image build run run-container clean push-image test

bin := ./bin/roomloggo
image_built := ./bin/image_built
image_pushed := ./bin/image_pushed
go_files := $(shell find . -name "*.go")
image_tag := roomloggo

build: $(bin)


run: $(bin)
	$(bin)


image: $(image_built)


run-container: $(image_built)
	docker run -v $(shell pwd)/roomloggo.config.yaml:/app/roomloggo.config.yaml $(image_tag)


clean:
	rm -rf ./bin


push-image: $(image_pushed)


test:
	go test ./...


$(image_pushed): $(image_built)
	docker save $(image_tag) | pv | ssh raspi docker load
	touch $(image_pushed)

$(bin): $(go_files)
	CGO_ENABLED=1 go build -o $(bin)


$(image_built): $(go_files)
	docker buildx build --platform linux/arm64/v8 --tag $(image_tag) .
	touch $(image_built)