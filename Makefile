all: gen test install

help:
	@echo "\tdeps: installs all dependencies"
	@echo "\tgen: generates boilerplate code"
	@echo "\ttest: Run all tests"

deps:
	@echo "---------------------------"
	@echo "Performing dependency check"
	@echo "---------------------------"
	go get -u golang.org/x/tools/cmd/goimports
	go get -u -t ./...                                  # get the application with all its deps

generate:
	@echo "----------------------"
	@echo "Generating source-code"
	@echo "----------------------"
	go generate ./...

imports:
	@echo "------------------"
	@echo "Optimizing imports"
	@echo "------------------"
	find . -name '*.go' -exec goimports -l -w -local github.com/ {} \;

format: imports
	@echo "----------------------"
	@echo "Formatting source-code"
	@echo "----------------------"
	find . -name '*.go' -exec gofmt -l -s -w {} \;

gen: generate imports format

check:
	@echo "---------------------"
	@echo "Perform static analysis"
	@echo "---------------------"
	go vet ./...

test: check
	@echo "---------------------"
	@echo "Running backend tests"
	@echo "---------------------"
	go test ./...                        # run unit tests
	make format

citest:
	@echo "---------------------"
	@echo "Running backend tests"
	@echo "---------------------"
	go get -u golang.org/x/tools/cmd/goimports
	go generate -tags ci  ./...
	make imports
	go test -tags ci ./...                        # run unit tests
	make format

coverage:
	@echo "----------------"
	@echo "Running coverage"
	@echo "----------------"
	./scripts/coverage.sh --html

install:
	@echo "----------------------------"
	@echo "Installing"
	@echo "----------------------------"
	go install ./...

.PHONY:
	help deps gen check test citest coverage install all
