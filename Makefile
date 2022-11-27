# Copyright shilixingchen@outlook Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   -  image        - build server image
#   -  clean_image  - clean server image
#   -  bin          - build server in /build/build_bin

# Tool commands (overridable)
GO_CMD             ?= go
DOCKER_CMD         ?= docker
DOCKER_COMPOSE_CMD ?= docker-compose

TAG=v2.0

# build server image
.PHONY: image
image:
	cd build && ./build.sh image -t $(TAG)

# clean server image
.PHONY: clean_image
clean_image:
	docker rmi -f shilixingchen/game:latest
	docker rmi -f shilixingchen/shilixingchen:$(TAG)

# build server bin
.PHONY: bin
bin:
	cd build && ./build.sh bin
