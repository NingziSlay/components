# components

## install

```shell
go get -u github.com/NingziSlay/components
```

## config

usage:

```golang
package main

import (
	"github.com/NingziSlay/components"
	"log"
)

type Config struct{}

func main() {
	var c Config
	if err := components.MustMapConfig(&c); err != nil {
		log.Fatal(err)
	}
	// use c here
}
```