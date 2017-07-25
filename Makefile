EXAMPLES_DIR = examples/
ANNOTATIONS_DIR = annotations/

.PHONY: all annotations examples clean

all: annotations examples

annotations:
	@make -C $(ANNOTATIONS_DIR)

examples:
	@make -C $(EXAMPLES_DIR)

clean:
	@make clean -C $(ANNOTATIONS_DIR)
	@make clean -C $(EXAMPLES_DIR)

re: clean all
