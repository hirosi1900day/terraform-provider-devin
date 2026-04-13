default: build

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status ./...

test:
	go test -v ./...

build:
	go build -o terraform-provider-devin

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hirosi1900day/devin/1.0.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-devin ~/.terraform.d/plugins/registry.terraform.io/hirosi1900day/devin/1.0.0/$$(go env GOOS)_$$(go env GOARCH)/

clean:
	rm -f terraform-provider-devin

.PHONY: fmt vet lint test build install clean
