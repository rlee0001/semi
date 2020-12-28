package scopes

import (
    "github.com/llir/llvm/ir"
    "github.com/llir/llvm/ir/value"

    "semi/src/types"
)

// A combination Name/Type used as the key in scopes.
type Declaration interface {
    Name() string
    Type() types.Type
}

// Object represents any name that can exist in a scope
type Object interface {
    Type() types.Type
}

// Object backed by an llvm ir Func
type Function interface {
    Object

    IrFunc() *ir.Func
}

// Object backed by an llvm ir Value
type Value interface {
    Object

    IrValue() value.Value
}

// Lookup(search-signature)
//  0. Let found-signature equal nil
//  1. Let candidate-scope equal the lexical-scope.
//  2. For each candidate-signature in the candidate-scope that matches the name of the search-signature:
//     a. Let candidate-signature equal the signature corresponding to the name found.
//     b. If the type-arity of the candidate-signature does not equals that of the search-signature:
//          i. Emit a signature type-arity error.
//     c. If any of the types of the candidate-signature are un-assignable to the types of the search-signature:
//          i. Go to the next candidate-signature at (2).
//     d. If the candidate-signature's types are more precise than the found-signature types:
//          i. Let the found-signature equal to the candidate-signature.
//  3. If the candidate-scope is not the root scope:
//      i. Let the candidate-scope equal the parent of the candidate-scope.
//      ii. Go to (2).
//  4. If the found-signature is nil:
//      i. Emit a signature error.
//  5. Return the Object associated with the found-signature.

// Integer v = 0;

// T w<T Integer> = 1;
// T w<T Float> = 3.14;
// fun goo(x Integer) {};
// goo w;
// goo w<Integer>;

// interface Map<K, V> { fun get(K) V; fun put(k, V); fun has(K) Boolean; };
// Map<String, T> x<T Integer> = ...;
// Map<String, T> x<T Float> = ...;
// fun bar(foo String) Boolean { return x<Integer>.has(foo); };
// fun baz(foo String) Integer { return x.get(foo); };
// fun buzz(foo String, fu Float) { x.put(foo, fu); };

// interface Function<P (), R> { fun invoke P R; };
// Function<(String), T> y<T Integer> = ...;
// Function<(String), T> y<T Float> = ...;
// fun y(foo String) String { return "" };              // syntax sugar for: Function<P, R> y<P (String), R String> = ...;
// Integer z = y("");
// Float z = y("");
// String z = y("");

// Scope represents a scope of named values such as functions, globals, constants, parameters and locals
type Scope interface {
    // Lookup looks up a name, returning its associated value if the name is found, otherwise returning an error
    Lookup(declaration Declaration) (Object, error)

    // DeclareValue creates a new value name in the scope
    DeclareValue(declaration Declaration, value Value) error

    // DeclareFunction creates a new function name in the scope
    DeclareFunction(declaration Declaration, function Function) error

    // AssignValue assigns a new value to the specified name
    AssignValue(declaration Declaration, value Value) error
}

// ScopeStack represents a stack of nested scopes such that lookups that fail to find a match in the top of the stack
// automatically cascade to the parent scope until the bottom of the stack is reached
type ScopeStack interface {
    Scope

    // Push pushes a new empty Scope on the stack
    Push()

    // Pop pups the top-most scope off the stack
    Pop()
}
