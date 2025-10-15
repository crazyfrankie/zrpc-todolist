.PHONY: buf-gen errcode

errcode:
	@echo "Generating error code"
	@./scripts/gen-error.sh --biz "*"

buf-gen:
	@echo "Generating buf code"
	@protoc --go_out=./protocol --go-zrpc_out=./protocol ./idl/auth.proto && protoc --go_out=./protocol --go-zrpc_out=./protocol ./idl/task.proto && protoc --go_out=./protocol --go-zrpc_out=./protocol ./idl/user.proto
