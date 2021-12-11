REGISTRY := ghcr.io/softree-group/protocall-connector

all: clerk connector

.PHONY: clerk
clerk:
	docker build \
		-f build/clerk.Dockerfile \
		-t $(REGISTRY)/clerk:$(VERSION) \
		--build-arg IMAGE=clerk \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

.PHONY: push-clerk
push-clerk: clerk
	docker push $(REGISTRY)/clerk:$(VERSION)

.PHONY: connector
connector:
	docker build \
		-f build/connector.Dockerfile \
		-t $(REGISTRY)/connector:$(VERSION) \
		--build-arg IMAGE=connector \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

.PHONY: push-connector
push-connector: connector
	docker push $(REGISTRY)/connector:$(VERSION)

.PHONY: push
push: push-connector push-clerk
