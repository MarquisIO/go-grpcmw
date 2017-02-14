EXAMPLES_DIR = examples/

.PHONY: examples

examples:
	@make -C $(EXAMPLES_DIR)

clean:
	@make clean -C $(EXAMPLES_DIR)
