PROTOC = protoc

PROTO_SRC = proto/example.proto
PROTO_PB_GO = $(PROTO_SRC:.proto=.pb.go)
PROTO_PB_MW_GO = $(PROTO_SRC:.proto=.pb.mw.go)

.PHONY: all clean re

all: $(PROTO_PB_GO) $(PROTO_PB_MW_GO)

clean:
	$(RM) $(PROTO_PB_GO) $(PROTO_PB_MW_GO)

re: clean all


%.pb.go: %.proto
	$(PROTOC) -I $(GOPATH)/src:. --go_out=plugins=grpc:. $^

%.pb.mw.go: %.proto
	$(PROTOC) -I $(GOPATH)/src:. --grpc-middleware_out=:. $^
