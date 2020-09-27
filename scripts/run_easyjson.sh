# !/bin/sh

go get -u github.com/mailru/easyjson/...
$(go env GOPATH)/bin/easyjson -all internal/api/response.go
$(go env GOPATH)/bin/easyjson -all internal/client/response.go
$(go env GOPATH)/bin/easyjson -all internal/gateway/gateways.go