.PHONY: image build run run-container clean push-image test

bin_dir := ./bin
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
	rm -rf $(bin_dir)


push-image: $(image_pushed)


test:
	go test ./...


$(bin_dir):
	mkdir -p $(bin_dir)


$(image_pushed): $(image_built)
	docker save $(image_tag) | pv | ssh raspi docker load
	touch $(image_pushed)


$(bin): $(go_files) | $(bin_dir)
	CGO_ENABLED=1 go build -o $(bin)


$(image_built): $(go_files) | $(bin_dir)
	docker buildx build --platform linux/arm/v7 --tag $(image_tag) .
	touch $(image_built)