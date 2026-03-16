BUILD_IMAGE := cuda-build:latest
RUN_IMAGE   := cuda-run:latest
BUILDER     := cuda-java-builder

.PHONY: all build run builder clean

all: builder

build:
	docker build -f build.Dockerfile -t $(BUILD_IMAGE) .

run:
	docker build -f run.Dockerfile -t $(RUN_IMAGE) .

builder: build run
	pack builder create $(BUILDER) --config builder.toml

clean:
	docker rmi $(BUILD_IMAGE) $(RUN_IMAGE) || true
	pack builder remove $(BUILDER) || true
