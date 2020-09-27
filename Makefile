.PHONY: all clean nocache run mocker run-mocker test test-integration

CWD=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CONF=$(CWD)/configs/config.yaml
SSLPATH=$(CWD)/configs/ssl
MOCKCONF=$(CWD)/configs/mockconfig.yaml

OUT=bin/payservice
MOCKER=bin/mocker

all: easyjson
	go build -v -o $(CWD)/$(OUT) $(CWD)/cmd/main.go

run: export PAYSERVICE_CONFIG_PATH=$(CONF)
run: export SSL_KEY=$(SSLPATH)/localhost.key
run: export SSL_CERT=$(SSLPATH)/localhost.crt
run: all
	$(CWD)/$(OUT)

mocker:
	go build -v -o $(CWD)/$(MOCKER) $(CWD)/test/mockery/main.go

run-mocker: export PAYSERVICE_MOCKCONFIG_PATH=$(MOCKCONF)
run-mocker: mocker
	$(CWD)/$(MOCKER)

clean: nocache
	find . -name "*easyjson*.go" -delete
	rm -rf $(CWD)/bin/

test: nocache
	sh $(CWD)/scripts/run_test_unit.sh

test-integration: export PAYSERVICE_CONFIG_PATH=$(CONF)
test-integration: export SSL_KEY=$(SSLPATH)/localhost.key
test-integration: export SSL_CERT=$(SSLPATH)/localhost.crt
test-integration: nocache
	sh $(CWD)/scripts/run_test_integration.sh

nocache:
	go clean -testcache

easyjson:
	sh $(CWD)/scripts/run_easyjson.sh