proto: ## Compile protobuf file to generate Go source code for gRPC Services
	protoc --go_out=plugins=grpc:. helloworld/*.proto
##  protoc --go_out=plugins=grpc:. *.proto

ID?=1
.DEFAULT_GOAL := help

.EXPORT_ALL_VARIABLES:
PORT?=50051
HOST?=localhost
CAFILE?="ca.cert"

.PHONY: proto

all: cert docker-build

install-docker: ## Install Docker
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	OS_CODENAME=$(shell lsb_release -cs)
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu ${OS_CODENAME} stable"
	sudo apt-get update
	sudo apt-get install -y docker-ce
	sudo usermod -aG docker ${USER}


ssl: ## Create certificates to encrypt the gRPC connection
#	rm -rf conf && mkdir SSL && cd SSL
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext
	clear

rmssl:
	rm ca.* service.*
	clear