GO=go
ANTLR=antlr
BINARY_NAME=semi

all: test

gen:
	@$(ANTLR) -Dlanguage=Go -o src/parser/gen/ -package gen -Xexact-output-dir src/parser/Semi.g4

vendor:
	@$(GO) mod vendor

build: gen vendor
	@$(GO) build -o $(BINARY_NAME) src/semi.go

test/unit: build
	@$(GO) test -coverprofile cp.out ./...

test/build: test/unit
	@./$(BINARY_NAME) build -o examples/simple_print_expression examples/simple_print_expression.semi

test: test/unit test/build
	@test "$(shell ./examples/simple_print_expression)" = "$(shell printf '%f\n' -12)" || { exit 2; }

install:
	@cp $(BINARY_NAME) /usr/local/bin/

clean:
	@$(GO) clean
	@rm -f $(BINARY_NAME)
	@rm -rf src/parser/gen/
	@rm -f examples/simple_print_expression

.PHONY: all gen vendor build test/unit test/build test install clean
