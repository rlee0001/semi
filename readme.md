### Prerequisites

Make sure that the following dependencies are installed and working prior to
attempting to build or run semi:

 Dependency | Version | Installation (Mac OS X/Homebrew)
 -----------|---------|---------------------------------
 Golang     | 1.15.1  | `brew install go`
 Antlr      | 4.8_1   | `brew install antlr`
 llvm       | 10.0.1  | `brew install llvm`

### Building

To build (and test) `semi`, just run:

```
$ make
```

To regenerate `antlr`-generated sources, run:

```
$ make gen
```

### Running

To use `semi` to compile a `.semi` source file as an executable, run:

```
$ ./semi build SOURCE
```

...where `SOURCE` is the path to a `.semi` source file.

### TODO

Things that need to be done or completed, by approximate priority and required effort.

#### High Priority

 Functional Area | Effort | Description
 ----------------|--------|--------------------------------------------------------------------------
 Everywhere      | SMALL  | Comments on at least all exports
 Parser/CodeGen  | MEDIUM | Function Declaration and Calls
 CodeGen         | MEDIUM | Data model for types (unify function and value types)
 CodeGen         | MEDIUM | Function parameters & arguments
 Everywhere      | LARGE  | Unit tests

#### Medium Priority

 Functional Area | Effort | Description
 ----------------|--------|--------------------------------------------------------------------------
 Ast             | MEDIUM | Processing passes (validation/type checks, type-specific node types (add-int vs add-float), type coercions, syntactic sugar expansions)
 Parser          | MEDIUM | Custom error handling (check antlr API)
 Parser          | MEDIUM | Grammar Fixes (call statement, assignment statement)
 CodeGen         | MEDIUM | Function inlining (needed for standard library; all blocks are functions)
 CodeGen         | MEDIUM | Semi library implementation (if, loops, types, etc...)
 Everywhere      | LARGE  | Features (named operators, strings, arrays, structs, comments, lambdas, clauses, blocks, generics)

#### Low Priority

 Functional Area | Effort | Description
 ----------------|--------|--------------------------------------------------------------------------
 Everywhere      | SMALL  | Add `make fmt` and `make lint`
 Sub-commands    | SMALL  | Implement `semi run ...`
 Sub-commands    | SMALL  | Implement `semi version`
 Sub-commands    | SMALL  | Ensure llvm/clang are present, and check their versions
 Sub-commands    | SMALL  | Finish parse-tree sub-command
 Parser          | SMALL  | Add `--` (STDIN) argument support
 Parser          | SMALL  | Add `-c CODE` option support
 Parser          | MEDIUM | Find a way to only walk the parse tree once (refactor listeners...again)
 Examples        | MEDIUM | Example programs (with makefile tests)
 CodeGen         | LARGE  | Heap allocations, garbage collection
 CodeGen         | LARGE  | Runtime implementation (other than libc, replace clang with llc/ld)?

### Types

Type specifier vs definition expression:

 - Specifier: names a type, with all requisite generic arguments.
    Syntax: `( "[" "]" )? <name> ( "[" <type-specifier> ( "," <type-specifier> )+ "]" )?`
 - Definition: defines a type as one of: `struct{...}`, `interface{...}`, `func(...) ...`

Type:
    Name() string
    Ordinal() uint32
    Parameters() []TypeParameter
