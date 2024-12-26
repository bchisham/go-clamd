go-clamd
========

Interface to clamd (clamav daemon). You can use go-clamd to implement virus detection capabilities to your application.

[![GoDoc](https://godoc.org/github.com/dutchcoders/go-clamd?status.svg)](https://godoc.org/github.com/dutchcoders/go-clamd)
[![Build Status](https://travis-ci.org/dutchcoders/go-clamd.svg?branch=master)](https://travis-ci.org/dutchcoders/go-clamd)

## Installation

```bash
go get github.com/bchisham/go-clamd
```

## New Features

Now supports context.Context for canceling scans. A context can be passed to the ContextScanStream function to cancel a scan.

## Examples

```go
package main

import (
	"bytes"
	"context"
	"log"
	"time"
	"fmt"
	"github.com/bchisham/go-clamd"
)

func main() {
	var (
		err error
		res *clamd.ScanResult
	)
	// Connect to clamd server through a unix socket
	c := clamd.NewClamd("/tmp/clamd.socket")
	// scan a stream of bytes
	reader := bytes.NewReader(clamd.EICAR)

	cancelContext, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer func() {
		cancelFunc()

	}()
	ss, err := c.ContextScanStream(cancelContext, reader)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case <-cancelContext.Done():
			err = fmt.Errorf("scan cancelled")
		case res = <-ss:
			if res == nil {
				log.Fatalln("nil result")
			}
		}
		switch res.Status {
		case clamd.RES_OK:
			log.Println("no virus found")
		case clamd.RES_FOUND:
			log.Printf("virus found: %s\n", res.Raw)
		case clamd.RES_ERROR, clamd.RES_PARSE_ERROR:
			log.Fatalf("error: %s\n", res.Raw)
		}
	}
	return
}
```

## Contributions

Contributions are welcome.

## Creators 

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>

- <https://twitter.com/dutchcoders>

## Contributors

**Brandon Chisham**
- <https://github.com/bchisham>
- <https://twitter.com/run4ever79>


## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef. Code released under [the MIT license](LICENSE). 
