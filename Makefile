GIT_DATE := $(shell git log -1 --pretty='%ci')
GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_VERSION := version: ${GIT_HASH}, date: $(GIT_DATE)

SFX_PLUGIN_CONFIG_PKG := github.com/signalfx/signalfx-fluent-bit-plugin/internal/config

binary:
	go build -buildmode=c-shared \
	-ldflags '-X "$(SFX_PLUGIN_CONFIG_PKG).gitVersion=$(GIT_VERSION)"' \
	-o dist/signalfx.so

test:
	go test ./... -cover

cover:
	mkdir -p dist
	go test -coverprofile=dist/c.out ./...
	go tool cover -html=dist/c.out -o dist/coverage.html
	open dist/coverage.html

images:
ifndef TAG
	$(error Docker tag not set. Syntax: make images TAG=x.y.z)
endif
	docker build . -t quay.io/signalfx/fluent-bit:$(TAG) -f ./build/package/Dockerfile
	docker build . -t quay.io/signalfx/fluent-bit-aws:$(TAG) -f ./build/package/Dockerfile.aws

publish:
ifndef TAG
	$(error Docker tag not set. Syntax: make publish TAG=x.y.z)
endif
	docker push quay.io/signalfx/fluent-bit:$(TAG)
	docker push quay.io/signalfx/fluent-bit-aws:$(TAG)

demo:
	docker run -it --rm quay.io/signalfx/fluent-bit:$(TAG) fluent-bit/bin/fluent-bit -i cpu -t cpu \
	-F record_modifier -p 'Whitelist_key=cpu_p' -m '*' \
	-e /fluent-bit/signalfx.so -o SignalFx -p IngestURL=$(INGEST_URL) -p Token=$(TOKEN) \
	-p MetricName=com.example.app.requests -m '*' -p LogLevel=debug

clean:
	rm -rf dist
	rm -f signalfx-fluent-bit-plugin
