LATEST=natw/debug-world:latest

build:
	docker build . -t $(LATEST)

RELEASE_TAG="natw/debug-world:$(TAG)"

release:
ifndef TAG
	$(error MISSING TAG: please specify a version to tag release with (ie `make release TAG=v123`))
endif
	docker tag $(LATEST) $(RELEASE_TAG)
	docker push $(RELEASE_TAG)

.PHONY: build release
