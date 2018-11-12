all: analyzer/analyzer apiserv/apiserv

test: analyzertest apiservtest

install:
	scripts/install.sh $(INSTALLFLAGS) --prefix $(PREFIX)

analyzer/analyzer: analyzer/data.go analyzer/static.go
	cd analyzer && go build

apiserv/apiserv: apiserv/main.go
	cd apiserv && go build

analyzertest:
	cd analyzer && go test

apiservtest:
	cd apiserv && go test
