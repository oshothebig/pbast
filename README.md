# pbast [![Build Status](https://travis-ci.org/oshothebig/pbast.svg?branch=master)](https://travis-ci.org/oshothebig/pbast)

Simple AST library for Protocol Buffers (proto3) in Golang

## Description
This package provides constructs defined [Protocol Buffers Version 3 Language Specification](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec).
It is designed to create a Protocol Buffers' AST by those constructs, but not designed to parse ".proto" files.
One of the typical use cases is builing a Protocol Buffers' AST when transforming an AST defined for a different language.
`printer` sub-package allows us to output an AST to `io.Writer` in Protocol Buffers' file format.

## Install
This package is "go gettable".

`go get github.com/oshothebig/pbast`

## Reference
- Protocol Buffers
    - [Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
    - [Protocol Buffers Version 3 Language Specification](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec)

## License
Apache License Version 2.0