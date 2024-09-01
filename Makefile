TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=xtratuscloud
NAMESPACE=local
NAME=azureipam
BINARY=terraform-provider-${NAME}
VERSION=1.1.0
OS_ARCH=linux_amd64

default: install

clean:
	rm -rf dist/${BINARY}
	
build: clean
	go build -o dist/${BINARY} -ldflags="-X 'main.Version=v${VERSION}'"

release:
	goreleaser release --clean --snapshot --skip=sign,publish

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m