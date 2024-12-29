proto-gen-server:
    @echo "Generating protobuf stubs for server"
    buf generate

proto-gen-python:
    @echo "Generating protobuf stubs for Python client"
    python -m grpc_tools.protoc \
           -I=proto \
           --grpc_python_out=clients/python \
           --pyi_out=clients/python \
           --python_out=clients/python \
           ./proto/evochi/v1/evochi.proto
    ruff format ./clients/python/evochi/v1/evochi_pb2*
