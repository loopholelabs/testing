# TestConn


[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/loopholelabs/testconn)](https://goreportcard.com/report/github.com/loopholelabs/testconn)
[![go-doc](https://godoc.org/github.com/loopholelabs/testconn?status.svg)](https://godoc.org/github.com/loopholelabs/testconn)

TestConn is a [Go](http://golang.org) library that creates a pair of TCP Connections (that satisfy the `net.Conn` interface) that are bound to one another. It is simply a helper package 
designed to be used in testing and cleans up the listener after itself while watching for race conditions. 

**This library requires Go 1.17 or later.**

## Important note about releases and stability

This repository generally follows [Semantic
Versioning](https://semver.org/). However, **this library is currently in Alpha** and
is still considered experimental. Breaking changes of the library will _not_ trigger a
new major release. The same is true for selected other new features explicitly marked as
**EXPERIMENTAL** in CHANGELOG.md.

## Usage

```go
package main 

import (
	"github.com/loopholelabs/testconn"
	"math/rand"
)

func main() {
	// Create a byte slice of random data
	data := make([]byte, 512)
	rand.Read(data)

	// Use the testconn.New() function to get a 
	// new pair of connections
	c1, c2, err := testconn.New()
	if err != nil {
		panic(err)
    }

	// c1 and c2 are real TCP connections that satisfy the `net.Conn` interface
	_, err = c1.Write(data)
	if err != nil {
		panic(err)
	}

	read := make([]byte, 512)
	_, err = c2.Read(read)
	if err != nil {
		panic(err)
	}
}
```

## License

The TestConn project is available as open source under the terms of the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).


## Project Managed By:
[![https://loopholelabs.io][LOOPHOLELABS]](https://loopholelabs.io)

[GITREPO]: https://github.com/loopholelabs/testconn
[LOOPHOLELABS]: https://cdn.loopholelabs.io/loopholelabs/LoopholeLabsLogo.svg
