# protobuf-stub-example

0. Install protoc
    a. https://grpc.io/docs/protoc-installation/
    b. `go install google.golang.org/protobuf/cmd/protoc-gen-go`
    c. `go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0`
1. Generate Go code from protobuf:

    ```
    cd /path/to/protobuf-stub-example
    protoc --go_out=plugins=grpc:. *.proto
    ```