.PHONY: image build run run-container clean push-image test

bin_dir := ./bin
bin := ./bin/roomloggo
image_built := ./.make/image_built
image_pushed := ./.make/image_pushed
make_dir := ./.make/
go_files := $(shell find . -name "*.go")
image_tag := roomloggo
target_platform := linux/arm64

build: $(bin)


run: $(bin)
	$(bin)


image: $(image_built)


run-container: $(image_built)
	podman run -v $(shell pwd)/roomloggo.config.yaml:/app/roomloggo.config.yaml $(image_tag)


clean:
	rm -rf $(bin_dir)
	rm -rf $(make_dir)


push-image: $(image_pushed)


test:
	go test ./...


$(bin_dir):
	mkdir -p $(bin_dir)


$(image_pushed): $(image_built)
	podman save $(image_tag) | pv | ssh raspi docker load
	touch $(image_pushed)


$(bin): $(go_files) | $(bin_dir)
	CGO_ENABLED=1 go build -o $(bin)


$(image_built): $(go_files) | $(bin_dir)
	podman build --platform $(target_platform) --tag $(image_tag) .
	touch $(image_built)
