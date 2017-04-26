EXAMPLES_DIR = examples/
PROTO_DIR = proto/

.PHONY: all proto examples clean

all: proto examples

proto:
	@make -C $(PROTO_DIR)

examples:
	@make -C $(EXAMPLES_DIR)

clean:
	@make clean -C $(PROTO_DIR)
	@make clean -C $(EXAMPLES_DIR)

re: clean all
