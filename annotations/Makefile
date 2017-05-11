PROTOC = protoc

PROTO_SRC = annotations.proto
PROTO_PB_GO = $(PROTO_SRC:.proto=.pb.go)

DESCRIPTOR = github.com/golang/protobuf/protoc-gen-go/descriptor

GO_PACKAGE = annotations

.PHONY: all clean re

all: $(PROTO_PB_GO)

clean:
	$(RM) $(PROTO_PB_GO)

re: clean all

%.pb.go: %.proto
	$(PROTOC) --go_out=import_path=$(GO_PACKAGE),Mgoogle/protobuf/descriptor.proto=$(DESCRIPTOR),plugins=grpc:. $^
