PROTOC = protoc

PROTO_SRC = interceptors.proto
PROTO_PB_GO = $(PROTO_SRC:.proto=.pb.go)

DESCRIPTOR = github.com/golang/protobuf/protoc-gen-go/descriptor

.PHONY: all clean re

all: $(PROTO_PB_GO)

clean:
	$(RM) $(PROTO_PB_GO)

re: clean all

%.pb.go: %.proto
	$(PROTOC) --go_out=Mgoogle/protobuf/descriptor.proto=$(DESCRIPTOR),plugins=grpc:. $^
