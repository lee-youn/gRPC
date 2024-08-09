PIP=pip3
PYTHON=python3

ibrew:
	/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

iprotoc:
	brew install protobuf

iprotoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ipyproto:
	$(PIP) install grpcio
	$(PIP) install grpcio-tools
	


all:
	make -C client_go protos
	make -C client_go_call_python protos
	make -C client_py protos
	make -C server_go protos
	make -C server_py protos
