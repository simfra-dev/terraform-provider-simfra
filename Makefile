default: build

build:
	go build -o terraform-provider-simfra

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/simfra-dev/simfra/0.0.1/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-simfra ~/.terraform.d/plugins/registry.terraform.io/simfra-dev/simfra/0.0.1/$$(go env GOOS)_$$(go env GOARCH)/

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

generate:
	go generate ./...

.PHONY: default build install test testacc generate
