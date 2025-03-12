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

proto-gen: proto-gen-server proto-gen-python

bump version:
    @echo "Bumping version to {{ version }}"
    echo {{ version }} > "./version"
    sed -i "s/^version = .*/version = \"$(echo {{ version }} | sed 's/^v//')\"/" ./clients/python/pyproject.toml
